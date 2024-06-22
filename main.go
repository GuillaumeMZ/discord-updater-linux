package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"
)

func loop_until_connected() {
	for {
		res, err := http.Head("https://www.google.com")

		if err == nil && res.StatusCode == 200 {
			break
		}

		time.Sleep(time.Second)
	}
}

func fetch_latest_discord_version() string {
	res, err := http.Head("https://discord.com/api/download?platform=linux&format=deb")

	if err != nil {
		panic("Error: couldn't fetch latest Discord version: HEAD request failed.")
	}

	//res.Request.URL is the actual download URL which contains the version
	regex, _ := regexp.Compile("https://dl.discordapp.net/apps/linux/(\\d+\\.\\d+\\.\\d+)/.*")
	matches := regex.FindStringSubmatch(res.Request.URL.String())

	if matches == nil {
		panic("Error: couldn't extract Discord's latest version from the download URL.")
	}

	return matches[1]
}

func parse_build_info(build_info_path string) map[string]interface{} {
	build_info, err := os.Open(build_info_path)
	defer build_info.Close()

	if err != nil {
		panic("Error: couldn't open build_info.json")
	}

	build_info_contents, err := io.ReadAll(build_info)

	if err != nil {
		panic("Error: couldn't read build_info.json")
	}

	var data map[string]interface{}
	json_err := json.Unmarshal(build_info_contents, &data)

	if json_err != nil {
		panic("Error: failed to parse build_info.json")
	}

	return data
}

func load_installed_discord_version(build_info_path string) string {
	data := parse_build_info(build_info_path)
	version := data["version"].(string)

	if version == "" {
		panic("Error: field \"version\" not found in build_info.json.")
	}

	return version
}

func update_local_discord_version(build_info_path string, new_version string) {
	data := parse_build_info(build_info_path)
	data["version"] = new_version

	build_info, err := os.OpenFile(build_info_path, os.O_WRONLY|os.O_CREATE, 0644)
	defer build_info.Close()

	if err != nil {
		panic("Error: couldn't create build_info.json.")
	}

	output_json, err := json.MarshalIndent(data, "", "\t")

	if err != nil {
		panic("Error: couldn't marshal the new build_info.json.")
	}

	build_info.Write(output_json)
}

func main() {
	if len(os.Args) != 2 {
		panic("Error: only one argument expected: the absolute path to your Discord installation's build_info.json.")
	}

	loop_until_connected()

	build_info_path := os.Args[1]
	local_discord_version := load_installed_discord_version(build_info_path)
	latest_discord_version := fetch_latest_discord_version()

	if local_discord_version == latest_discord_version {
		os.Exit(0)
	}

	update_local_discord_version(build_info_path, latest_discord_version)
}
