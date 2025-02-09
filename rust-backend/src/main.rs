use axum::{
    routing::{get, post},
    Router,
};

use rust_backend::{db, middleware::auth::*, types::AppState};

#[tokio::main]
async fn main() {
    let pool = db::ops::connect();

    let app_state = AppState { pool };

    // build our application with a route
    let app = Router::new()
        .route("/login", post(login_handler))
        .route("/register", post(register_handler))
        .with_state(app_state);

    // run our app with hyper, listening globally on port 3000
    let listener = tokio::net::TcpListener::bind("0.0.0.0:3000").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}