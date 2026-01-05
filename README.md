# tennisaboplaner-go

Tennisaboplaner is a script to schedule matches over a season. It's purpose is to give a schedule, if you are group of people that meet weekly for a tennis match and have only a number of limited courts, for instance indoor / in the winter season. It searches for a schedule in such a way that all people play the same amount, every player faces all others evenly and the matches are evenly distributed per player during the season. Additionally, you can provide prescheduling dates where player can't play. This will be taken into account, when it searches for possible schedules. It is a reimplementation of <https://github.com/Elladur/tennisaboplaner> in Go.

## Requirements

Currently, we don't have any system requirements.

## Usage

Adapt settings.json to your needs and execute `tennisaboplaner-go`. You can use the `--help` to see default parameters and possible options.
It will generate a Excel-File in the output-folder which represents the schedule. Additionally, it will create a calendar (*.ics) for each player.

## Contributing

Pull requests are welcome. Please open an issue first
to discuss what you would like to change.
It is possible to use the devcontainer configuration to clone into a container volume and build the source code. 

Please make sure to update tests as appropriate.
