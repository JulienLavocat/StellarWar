use galaxy::generator;
use spacetimedb::spacetimedb;

pub mod galaxy;

#[spacetimedb(reducer)]
pub fn generate_map() {
    generator::generate()
}
