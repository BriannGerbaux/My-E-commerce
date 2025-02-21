// @generated automatically by Diesel CLI.

diesel::table! {
    cart (id) {
        id -> Int4,
        user_id -> Nullable<Int4>,
        product_id -> Nullable<Int4>,
    }
}

diesel::table! {
    products (id) {
        id -> Int4,
        #[max_length = 20]
        name -> Varchar,
        description -> Text,
        price_in_dollar -> Int4,
        thumbnail_url -> Text,
        created_at -> Timestamp,
        updated_at -> Timestamp,
    }
}

diesel::table! {
    users (id) {
        id -> Int4,
        #[max_length = 20]
        username -> Varchar,
        #[max_length = 50]
        email -> Varchar,
        password_hash -> Text,
        created_at -> Timestamp,
        updated_at -> Timestamp,
    }
}

diesel::joinable!(cart -> products (product_id));
diesel::joinable!(cart -> users (user_id));

diesel::allow_tables_to_appear_in_same_query!(
    cart,
    products,
    users,
);
