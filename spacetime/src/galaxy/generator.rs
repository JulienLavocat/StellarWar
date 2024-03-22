use poisson_diskus::bridson_rng;
use rand::prelude::StdRng;
use rand_seeder::Seeder;

pub fn generate() {
    generate_systems()
}

fn generate_systems() {
    log::info!("generating systems");
    let box_size = [200.0, 200.0];
    let rmin = 50.0;
    let num_attempts = 30;

    let mut rng: StdRng = Seeder::from("test").make_rng();
    for point in bridson_rng(&box_size, rmin, num_attempts, false, &mut rng)
        .unwrap()
        .iter()
    {
        log::debug!("x: {}, y: {}", point[0], point[1])
    }
}
