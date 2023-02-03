use clap::Parser;
use simple_steam_totp::{generate};

fn find_default_steamcmd() -> &'static str {
    if cfg!(target_os = "windows") {
        "C:\\steamcmd\\steamcmd.exe"
    } else {
        if (std::path::Path::new("/home/steam/steamcmd/steamcmd.sh")).exists() {
            "/home/steam/steamcmd/steamcmd.sh"
        } else {
            "/home/steam/steamcmd"
        }
    }
}

#[derive(Parser, Debug)]
#[clap(version, about, long_about = None)]
struct Args {
    // Path to steamcmd executable
    #[clap(long, default_value = find_default_steamcmd())]
    path: String,

    // Steam username
    #[clap(short, long)]
    username: String,

    // Steam password
    #[clap(short, long)]
    password: String,

    // Steam 2FA shared secret
    #[clap(short, long)]
    secret: String,

    // Steamcmd args
    #[clap(short, long)]
    args: String,
}

fn main() {
    let args = Args::parse();

    if !std::path::Path::new(&args.path).exists() {
        println!("Steamcmd executable not found at {}. Please specify with --path", args.path);
        std::process::exit(1);
    }

    let totp = match generate(&args.secret) {
        Ok(code) => code,
        Err(e) => {
            println!("Failed to generate Steam TOTP code: {}", e);
            std::process::exit(1);
        }
    };

    let status = std::process::Command::new(&args.path)
        .arg("+login")
        .arg(&args.username)
        .arg(&args.password)
        .arg(&totp)
        .args(args.args.split(' '))
        .status();

    std::process::exit(status.unwrap().code().unwrap());
}