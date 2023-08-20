# Manaha Minder for Minecraft

Manaha is a server that I own/rent, the name is just a funny sound. 

In 2013 a friend convinced me to host a Minecraft instance, and this VPS was the easiest option.  
I came up with some scripts in shell and Perl in order to add some "nice" features - like having no operators by default, but allowing certain players to request temporary operator status.

In 2016 the server closed down due to lack of use.  Now in 2023 I started getting the itch to play again, and decided to resurrect the server.  The old scripts weren't very self-explanatory, so I decided to rewrite them in Go.

## Features

- Welcome message when players log in
- Optional random gift given when players log in
- No operators by default.  
  You can configure a list of players who are allowed to request operator status for a specified amount of time.
- Ability to specify Minecraft commands to run when players say something specific in chat.
- Website helpers:
  - Show online players.
  - For offline players, show when they were last seen.
  - Log activity to a file.
  - Calculate the amount of time each player has been online.
- Checks is the Minecraft server has crashed, and restarts it if necessary.
- Ability to restart the Minecraft server at a specified time each day.
- Unzips the historic log files from the Minecraft server.

## Installation

### Prerequisites

Manaha Minder was built to work with Minecraft servers managed with [MSM](https://msmhq.com/).  
Sending commands to the Minecraft server is done via MSM.

If you want to use the status and activity logging features, you will need a web server.

This has been designed for a Linux based system running systemd.  It may work on other systems, but this has not been tested.

### Installing

Currently, there are no binaries available, so you will need to build from source:

1. Note down the user that is running the Minecraft server via MSM.  This will be used later.
2. Clone the repository   
   ```shell
   git clone github.com/forquare/manaha-minder
   ```
3. Build the binary  
   ```shell
   go build
   ```
4. Copy the binary to the location you want to run it from  
   ```shell
   sudo cp manaha_minder /usr/local/bin
   ```
5. Create the configuration file "`/etc/mminder.yaml`", and use the configuration section below to populate it
6. Replace the default user in the systemd service file with the user you noted down in step 1:  
   ```shell
   sed -i 's/manaha-minecrafter/YOUR_USER/' manaha-minder.service
   ```
7. Copy the systemd service file  
   ```shell
   sudo cp manaha-minder.service /etc/systemd/system
   ```
8. Enable the service:  
   ```shell
   sudo systemctl enable manaha-minder.service
   ```
9. Start the service:  
   ```shell
   sudo systemctl start manaha-minder.service
   ```
10. Check the status of the service:  
   ```shell
   sudo systemctl status manaha-minder.service
   ```

## Configuration

The configuration file is in the [YAML format](https://en.wikipedia.org/wiki/YAML).  

Although not a requirement, it is advisable to use the `yamllint` command to check the syntax of the file to prevent any errors.  To do so, install the `yamllint` package, and then run the command:  
```shell
yamllint /etc/mminder.yaml
```

In the configuration below, the server name "manaha" is used as an example.  Please replace this with the name of your server.

```yaml
---
minecraft_server:
  # The server name in MSM
  server_name: manaha

  # Location of the latest.log file
  latest_log_file: /opt/msm/servers/manaha/logs/latest.log

  # Location of all logs
  log_dir: /opt/msm/servers/manaha/logs

  # Path to the world directory
  world_dir: /opt/msm/servers/manaha/worldstorage/world

  # Path to Minecraft whitelist.json
  whitelist_file: /opt/msm/servers/manaha/whitelist.json

  # Path to MSM command
  msm_binary: /usr/local/bin/msm

  # Settings for restarting the Minecraft server
  restart:
    # Whether to restart the server
    enabled: true

    # The time to restart the server using the cron format
    # https://en.wikipedia.org/wiki/Cron
    # This example will restart the server at 5am every day
    cron: "0 5 * * *"

  # Keep an eye on common server errors and restart the server if they occur
  watchdog: true

  # Decompress archived log files
  log_decompress: false

manaha_minder:
  # SQLite database file
  # The database is used to store a variety of information
  # including player activity and operator requests
  database: /home/manaha-minecrafter/var/manaha-minder.db
      
  # Set the logging level for manaha minder.
  # The levels are:
      # Trace - shows everything, will likely overwhelm you!
      # Debug - gives a good idea about what is going on
      # Info - shows noteworthy events
      # Warn - shows events that are not necessarily errors, but may be worth investigating
      # Error - shows errors that are not fatal
      # Fatal - shows fatal errors, this also crashes the service
      # Panic - similar to fatal, but worse
  # It is recommended to use info or warn for normal operation
  log_level: Warn

login:
  # The message to display when a player logs in
  # <PLAYER> will be replaced with the player's name
  welcome_message: |-
    Welcome <PLAYER>!
    We hope you have fun on the server!
    
  # Enable the server to give a random gift to players when they log in.
  # There is about a 3% chance of a gift being given.
  give_random_gift: true


operator:
  # Duration that operator status lasts for in seconds
  duration: 600  # 10 minutes
  
  # Players that can become operators
  players:
    - username1
    - other_username

activity:
  # Whether to log player activity, if false, the rest of the activity settings are ignored
  log_activity: true
  
  # Whether to recalculate the activity of players when the server starts
  recalculate_activity_on_startup: true
  
  ###################################################################
  # Website helpers                                                 #
  # These are optional, and require a web server to be configured   #
  # Please see the wiki for more information:                       #
  # https://github.com/forquare/Manaha-Minder/wiki/Webpage          #
  ###################################################################
  
  # Enable calculating the amount of time each player has been online and outputting it to a file
  generate_time_played_output: true
  
  # Where the activity calculations should be outputted to
  time_played_file: /var/www/html/activity.txt
  
  # Enable generating a status page for the website
  generate_status_table: true
  
  # Where the status page should be outputted to
  status_file: /var/www/html/players.html

# Custom actions are a YAML list of actions to perform when a command is run
custom_actions:
    # The name of the action
  - name: "Go home"
    # A description of the action
    description: "Teleport a player back to the original spawn point"
    # What the action should be triggered on - this uses a regular expression: https://en.wikipedia.org/wiki/Regular_expression
    pattern: "^take me home$"
    # The commands to run when the action is triggered
    # <PLAYER> will be replaced with the player's name
    # There is a 1 second delay between each command
    commands:
      - "say <PLAYER> wants to go home"
      - "say Teleporting PLAYER to original spawn point"
      - "tp <PLAYER> 0 73 0"

    # Another example
  - name: "Need some coal"
    description: "Give a the player some coal"
    pattern: "^need some coal$"
    commands:
      - "say <PLAYER> wants some coal"
      - "say Giving <PLAYER> some coal"
      - "give <PLAYER> coal"
```

## Running

Once you have installed and configured Manaha Minder, you can start it by running the command:  
```shell
sudo systemctl start manaha-minder.service
```

## Contributing

Please read [CONTRIBUTING.md](.github/CONTRIBUTING.md) for details on the code of conduct, and the process for submitting Pull Requests.

## Versioning

[Semantic Versioning](http://semver.org/) is used for versioning.

## Authors

- **Ben Lavery-Griffiths** - [Forquare](https://github.com/forquare) - [Website](https://hashbang0.com)

See also the list of [contributors](https://github.com/forquare/zht/contributors) who have participated in this project.

## License

This project is licensed under the [MIT Licence](LICENSE) - see the [LICENSE](LICENSE) file for details

## Releasing

1. Apply code changes via Pull Requests
2. Update local main branch to latest
3. Create an annotated tag with the new version  
   ```git tag -a x.x.z -m "Descriptive message"```
4. Push the tag  
   ```git push --tags```
5. The release action will run using GoReleaser to create the next release

All tags are [protected](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/managing-repository-settings/configuring-tag-protection-rules) so that only certain users can create them.
