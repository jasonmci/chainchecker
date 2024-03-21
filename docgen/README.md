# LanePairs Markdown Table Generator

This tool generates a markdown table that summarizes lane pairs and their corresponding source and destination based on a given configuration. It utilizes the lane pairs defined in `config.toml` or inputs provided by the user during runtime.

## Configuration

Before running the script, you should set up your lane pairs in the config.toml file. The format for the lane pairs is as follows:

```toml
LanePairs = [
    "Ethereum-Optimism",
    "Ethereum-Avalanche",
    ...
]
```
Uncomment or add new lane pairs as needed. The script will use these as default lane pairs if no input is provided during execution.

## Usage

1. avigate to the directory containing the script and the config.toml file.
1. Run the script using the following command:

```bash
go run . 
```

3. Upon execution, the script will prompt you to enter lane pairs. You can enter them as comma-separated values like LaneA-LaneB,LaneX-LaneY. If you wish to use the default lane pairs defined in config.toml, simply press Enter without typing anything.

```
Enter lane pairs (comma-separated, e.g., 'LaneA-LaneB,LaneX-LaneY'):
```

## Output
The markdown table includes the following headers:

* Lane Combo: The combination of the source and destination as defined by the user or the config file.
* Source: The starting point of the lane pair.
* Destination: The ending point of the lane pair.
* Additional headers "Scenario", "Status", and "Transactions" are included in the output for potential future enhancements but are not populated by the current version of the script.