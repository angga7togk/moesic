#!/usr/bin/env node

import { cac } from "cac";
import { PassThrough } from "stream";
import Speaker from "speaker";
import ffmpeg from "fluent-ffmpeg";
import boxen from "boxen";
import readline from "readline";
import {
	GITHUB_BASE_URL,
	MOESIC_V1_URL,
	RAW_GITHUB_BASE_URL,
} from "./constants.js";
import api from "./utils/api.js";
import utils from "./utils/index.js";
import { TextFormat } from "./utils/textformat.js";
import open from "open";

console.clear();

const cli = cac("moesic");

function makeControlsLine() {
	const left = `${TextFormat.BOLD}S${TextFormat.GRAY}kip${TextFormat.RESET}`;
	const centerLeft = `${TextFormat.BOLD}P${TextFormat.GRAY}ause${TextFormat.RESET}`;
	const centerRight = `${TextFormat.GRAY}S${TextFormat.RESET}${TextFormat.BOLD}o${TextFormat.GRAY}urce${TextFormat.RESET}`;
	const right = `${TextFormat.BOLD}Q${TextFormat.GRAY}uit${TextFormat.RESET}`;

	const totalWidth = 25;

	// Panjang real tanpa ANSI escape code
	const lengthLeft = utils.stripAnsi(left).length;
	const lengthCenterLeft = utils.stripAnsi(centerLeft).length;
	const lengthCenterRight = utils.stripAnsi(centerRight).length;
	const lengthRight = utils.stripAnsi(right).length;

	const totalContentLength =
		lengthLeft + lengthCenterLeft + lengthCenterRight + lengthRight;
	const spaceBetween = totalWidth - totalContentLength;

	const space1 = Math.floor(spaceBetween / 3);
	const space2 = Math.floor(spaceBetween / 3);
	const space3 = spaceBetween - space1 - space2;

	if (spaceBetween < 0) {
		return `${left} ${centerLeft} ${centerRight} ${right}`; // fallback
	}

	return (
		left +
		" ".repeat(space1) +
		centerLeft +
		" ".repeat(space2) +
		centerRight +
		" ".repeat(space3) +
		right
	);
}

function renderBox(title, elapsed, total) {
	const percent = Math.min((elapsed / total) * 100, 100);
	const barLength = 20;
	const filled = Math.floor((percent / 100) * barLength);
	const empty = barLength - filled;
	const bar = `[${"█".repeat(filled)}${"░".repeat(empty)}] ${
		TextFormat.GRAY
	}${utils.formatTime(elapsed)}/${utils.formatTime(total)}${TextFormat.RESET}`;

	return boxen(
		`${title}\n\n${bar}\n\n${makeControlsLine()}`,
		{
			padding: 1,
			width: 45,
			borderStyle: "singleDouble",
			borderColor: "cyanBright",
		}
	);
}

function getDurationFromUrl(url) {
	return new Promise((resolve, reject) => {
		ffmpeg.ffprobe(url, (err, metadata) => {
			if (err) return reject(err);
			resolve(metadata.format.duration);
		});
	});
}

export async function playSong(playlist, song) {
	const duration = Math.floor(await getDurationFromUrl(song.url));

	const speaker = new Speaker({
		channels: 2,
		bitDepth: 16,
		sampleRate: 44100,
	});

	const passthrough = new PassThrough();

	ffmpeg(song.url)
		.audioChannels(2)
		.audioFrequency(44100)
		.format("s16le")
		.on("error", (err) => {
			console.error("FFmpeg error:", err.message);
		})
		.pipe(passthrough);

	passthrough.pipe(speaker);

	let elapsed = 0;
	const interval = setInterval(async () => {
		readline.cursorTo(process.stdout, 0, 0);
		readline.clearScreenDown(process.stdout);
		console.log(renderBox(song.title, elapsed, duration));
		elapsed++;
		if (elapsed >= duration) {
			clearInterval(interval);
			const data = await api.getMoesicData(
				`${RAW_GITHUB_BASE_URL}/${MOESIC_V1_URL}`
			);
			const playlist = data[Math.floor(Math.random() * data.length)];

			const song =
				playlist.songs[Math.floor(Math.random() * playlist.songs.length)];
			passthrough.end();
			speaker.end();
			await playSong(playlist, song);
			clearInterval(interval);
		}
	}, 1000);

	process.stdin.setRawMode(true);
	process.stdin.resume();
	process.stdin.on("data", async (key) => {
		const input = key.toString().toLowerCase();
		if (input === "q") {
			clearInterval(interval);
			process.exit(0);
		} else if (input === "s") {
			const data = await api.getMoesicData(
				`${RAW_GITHUB_BASE_URL}/${MOESIC_V1_URL}`
			);
			const playlist = data[Math.floor(Math.random() * data.length)];

			const song =
				playlist.songs[Math.floor(Math.random() * playlist.songs.length)];
			passthrough.end();
			speaker.end();
			await playSong(playlist, song);
			clearInterval(interval);
		} else if (input === "p") {
			passthrough.pause();
		} else if (input === "o") {
			open(
				`${GITHUB_BASE_URL}/${MOESIC_V1_URL}#${utils.toSlug(
					playlist.title
				)}:~:text=${encodeURIComponent(song.title)}`
			);
		}
	});
}

cli.command("").action(async () => {
	const data = await api.getMoesicData(
		`${RAW_GITHUB_BASE_URL}/${MOESIC_V1_URL}`
	);
	const playlist = data[Math.floor(Math.random() * data.length)];

	const song =
		playlist.songs[Math.floor(Math.random() * playlist.songs.length)];
	await playSong(playlist, song);
});

cli.help();
cli.parse();
