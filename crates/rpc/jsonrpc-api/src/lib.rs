mod live_object;

pub mod server {
    pub use crate::live_object::LiveObjectApiServer;
}

pub mod client {
    pub use crate::live_object::LiveObjectApiClient;
}
