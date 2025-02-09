use axum::http::StatusCode;
use axum::response::{IntoResponse, Response};
use axum::Json;
use serde::Serialize;
use thiserror::Error;
use std::io;

use diesel::r2d2::{ConnectionManager, Pool};
use diesel::PgConnection;

#[derive(Clone)]
pub struct AppState {
    pub pool: Pool<ConnectionManager<PgConnection>>,
}

#[derive(Serialize)]
pub struct BasicResponse {
    message: String,
}

#[derive(Serialize)]
struct ErrorResponse {
    message: String,
}

#[derive(Debug, Error)]
pub enum Error {
    #[error("IO Error: {0}")]
    Io(#[from] io::Error),

    #[error("JWT Error: {0}")]
    JwtError(#[from] jsonwebtoken::errors::Error),

    #[error("Internal Error")]
    InternalError,

    #[error("Wrong password or email error")]
    WrongPasswordOrEmail,

    #[error("Email already in use")]
    EmailAlreadyInUse,

    #[error("Diesel error: {0}")]
    DieselError(#[from] diesel::result::Error),


}

impl IntoResponse for Error {
    fn into_response(self) -> Response {
        let (status, message) = match self {
            Error::Io(_) => (StatusCode::INTERNAL_SERVER_ERROR, "Internal server error"),
            Error::JwtError(_) => (StatusCode::INTERNAL_SERVER_ERROR, "Internal server error"),
            Error::InternalError => (StatusCode::INTERNAL_SERVER_ERROR, "Internal server error"),
            Error::WrongPasswordOrEmail => (StatusCode::BAD_REQUEST, "Wrong password or email"),
            Error::DieselError(_) => (StatusCode::BAD_REQUEST, "Wrong password or email"),
            Error::EmailAlreadyInUse => (StatusCode::BAD_REQUEST, "Email is already in use"),
        };

        let body = Json(ErrorResponse {
            message: message.to_string(),
        });


        (status, body).into_response()
    }
}