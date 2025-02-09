use crate::{db::models::{NewUser, Users}, types::{AppState, Error}};

use axum::{
    body::Body, extract::State, http::Response, Form
};
use argon2::{
    password_hash::{
        rand_core::OsRng, PasswordHash, PasswordHasher, PasswordVerifier, SaltString
    },
    Argon2
};
use diesel::prelude::*;
use serde::Deserialize;
use token::generate_jwt;

fn hash_password(text_password: &str) -> Result<String, Error> {
    let salt = SaltString::generate(&mut OsRng);

    // Argon2 with default params (Argon2id v19)
    let argon2 = Argon2::default();
    
    // Hash password to PHC string ($argon2id$v=19$...)
    let password_hash = argon2.hash_password(text_password.as_bytes(), &salt).map_err(|_| Error::InternalError)?.to_string();
    Ok(password_hash)
}

fn verify_password(text_password: &str, password_hash: &str) -> bool {
    let hash = PasswordHash::new(&password_hash).unwrap();
    Argon2::default().verify_password(text_password.as_bytes(), &hash).is_ok()
}

// /login body
#[derive(Deserialize)]
pub struct LoginBody {
    email: String,
    password: String
}

pub async fn login_handler(
    State(state): State<AppState>,
    Form(payload): Form<LoginBody>
) -> Result<String, Error> {
    use crate::db::schema::users::dsl::*;

    let pool = state.pool.clone();
    let result = tokio::task::spawn_blocking(move || {
        let mut conn = pool.get().map_err(|_| Error::InternalError)?;

        let user = users
            .filter(email.eq(payload.email))
            .limit(1)
            .select(Users::as_select())
            .load(&mut conn);

        match user {
            Ok(matched_user) => {
                if matched_user.len() > 0 && verify_password(&payload.password, &matched_user[0].password_hash) {
                    let jwt = generate_jwt(matched_user[0].id.to_string().as_str())?;
                    return Ok(format!("Bearer {}", jwt));
                } else {
                    return Err(Error::WrongPasswordOrEmail);
                }
            },
            Err(_) => return Err(Error::InternalError)

        }
    }).await.map_err(|_| Error::InternalError)?;
    result
}


// /register body
#[derive(Deserialize)]
pub struct RegisterBody {
    username: String,
    email: String,
    password: String
}

pub async fn register_handler(
    State(state): State<AppState>,
    Form(payload): Form<RegisterBody>
) -> Result<String, Error> {
    use crate::db::schema::users::dsl::*;

    let pool = state.pool.clone();
    let result = tokio::task::spawn_blocking(move || {
        let mut conn = pool.get().expect("Failed to get connection");

        let user = users
            .filter(email.eq(&payload.email))
            .limit(1)
            .select(Users::as_select())
            .load(&mut conn)
            .map_err(|e| Error::DieselError(e))?;

        if user.len() > 0 {
            return Err(Error::EmailAlreadyInUse);
        }

        let hash = hash_password(&payload.password)?;

        let new_user = NewUser { username: payload.username, email: payload.email, password_hash: hash };

        use crate::db::schema::users;
        diesel::insert_into(users::table)
            .values(&new_user)
            .returning(Users::as_returning())
            .get_result(&mut conn)
            .map_err(|e| Error::DieselError(e))?;

        Ok("User registered".to_string())
    })
    .await.map_err(|_| Error::InternalError)?;
    result
}

pub mod token {
    use serde::{Serialize, Deserialize};
    use jsonwebtoken::{encode, decode, Header, Algorithm, Validation, EncodingKey, DecodingKey};
    use chrono::{Utc, Duration};
    use std::env;

    use crate::types::Error;

#[derive(Debug, Serialize, Deserialize)]
pub struct Claims {
    sub: String,        // Subject (user ID)
    exp: usize,         // Expiration time (Unix timestamp)
}

/// Generates a JWT token
pub fn generate_jwt(user_id: &str) -> Result<String, Error> {
    let secret = env::var("JWT_SECRET").expect("JWT_SECRET must be set");

    let expiration = Utc::now()
        .checked_add_signed(Duration::minutes(60))  // Token valid for 60 minutes
        .expect("valid timestamp")
        .timestamp() as usize;

    let claims = Claims { sub: user_id.to_owned(), exp: expiration };

    let token = encode(
        &Header::default(), 
        &claims, 
        &EncodingKey::from_secret(secret.as_ref())
    ).map_err(Error::from)?;

    Ok(token)
}

/// Verifies and decodes a JWT token
pub fn verify_jwt(token: &str) -> Result<Claims, Error> {
    let secret = env::var("JWT_SECRET").expect("JWT_SECRET must be set");

    let decoded = decode::<Claims>(
        token,
        &DecodingKey::from_secret(secret.as_ref()),
        &Validation::new(Algorithm::HS256)
    ).map_err(Error::from)?;

    Ok(decoded.claims)
}
}