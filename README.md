# steamcmd-2fa

A simple tool that automatically generates a 2FA Steam Guard code and then runs `steamcmd` with that login and the rest of your arguments.

I wanted to run certain steamcmd commands during CI and did not want to face the hurdles of disabling Steam Guard, nor did I want to build an SMTP client that fetches the code, so I built this simple tool based on [go-steam-totp](https://github.com/fortis/go-steam-totp) that generates a code on-the-fly and runs steamcmd with it.

## Installation

You can either clone this and build it yourself or use the supplied binaries for both Windows and Linux.

```bash
git clone https://github.com/weilbyte/steamcmd-2fa
cd steamcmd-2fa
go build .
```

## Usage
```
Usage of steamcmd-2fa:
  -path string
        Path to steamcmd executable (default "C:\\steamcmd\\steamcmd.exe", "/home/steam/steamcmd", "/home/steam/steamcmd.sh")
  -code-only
        Only prints out the code without wrapping around steamcmd
  -username string
        Username to log in with
  -password string
        Password to log in with
  -seed string
        The 2FA seed/shared secret
  -args string
        Arguments to pass to steamcmd
```

For example, instead of running `steamcmd +login exampleuser examplepass +quit`, you would run `steamcmd-2fa --path /home/steam/steamcmd --username exampleuser --password examplepass --seed YOUR2FASEED --args "+quit"`. 

Or you can simply run `steamcmd-2fa --username exampleuser --password examplepass --seed YOUR2FASEED --code-only` to get only the code, which then you can supply to steamcmd's login command directly.

You can get your 2FA seed by [various methods](https://github.com/SteamTimeIdler/stidler/wiki/Getting-your-%27shared_secret%27-code-for-use-with-Auto-Restarter-on-Mobile-Authentication). Your seed here is the `shared_secret`.

**Please keep your seed/secret VERY safe as people with your seed can easily bypass your 2FA by simply generating the Steam Guard codes themselves!**
*I am not responsible if you somehow manage to lose/get yourself locked out of your account.*

## Contributing
Pull requests are welcome. 
## License
[MIT](https://choosealicense.com/licenses/mit/)
