use std::process;

fn main() {
    let args: Vec<String> = std::env::args().collect();
    if args.len() < 2 {
        eprintln!("usage: w3grs <replay.w3g>");
        process::exit(1);
    }
    let result = w3grs::parse_file(&args[1]).unwrap_or_else(|| {
        eprintln!("error: failed to parse replay");
        process::exit(1);
    });
    match serde_json::to_string_pretty(&result) {
        Err(e) => {
            eprintln!("error encoding JSON: {}", e);
            process::exit(1);
        }
        Ok(json) => println!("{}", json),
    }
}
