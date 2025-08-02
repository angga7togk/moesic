import React, {useState} from 'react';
import {Box, Text, useInput} from 'ink';
import crawler from './utils/crawler.js';

import {ChildProcessByStdio, spawn} from 'child_process';
import Stream from 'stream';
console.clear();

const items = [
	{
		title: 'Honey Sweet Tea Time - Tsumugi Kotobuki (Minako Kotobuki)',
		coverUrl:
			'https://i.ytimg.com/vi/tpDMYQCdiYA/hqdefault.jpg?sqp=-oaymwEnCOADEI4CSFryq4qpAxkIARUAAIhCGAHYAQHiAQoIGBACGAY4AUAB&rs=AOn4CLB34p9g5HwvSCJ9y2D2qNaHgI8Onw',
		provider: 'YOUTUBE',
		providerValue: 'tpDMYQCdiYA',
		malId: 5680,
	},
	{
		title: 'K on Azusa&Yui Fude pen Boru pen',
		coverUrl:
			'https://i.ytimg.com/vi/PCmZ45INWYA/hq720.jpg?sqp=-oaymwFBCNAFEJQDSFryq4qpAzMIARUAAIhCGAHYAQHiAQoIGBACGAY4AUAB8AEB-AH8CYAC0AWKAgwIABABGH8gJCgxMA8=&rs=AOn4CLAlc3vvNr0Y6Fxc7GdcuyBUP-MjuQ',
		provider: 'YOUTUBE',
		providerValue: 'PCmZ45INWYA',
		malId: 5680,
	},
	{
		title: 'K-On! | U&I - Hōkago Tea Time',
		coverUrl:
			'https://i.ytimg.com/vi/49WKHkLsR_U/hq720.jpg?sqp=-oaymwEnCNAFEJQDSFryq4qpAxkIARUAAIhCGAHYAQHiAQoIGBACGAY4AUAB&rs=AOn4CLBtPGgOSjeQ8LRhSbqNjNSptBonSQ',
		provider: 'YOUTUBE',
		providerValue: '49WKHkLsR_U',
		malId: 5680,
	},
	{
		title: 'K-On | Fuwa Fuwa Time',
		coverUrl:
			'https://i.ytimg.com/vi/I0xRbWqIohQ/hq720.jpg?sqp=-oaymwEnCNAFEJQDSFryq4qpAxkIARUAAIhCGAHYAQHiAQoIGBACGAY4AUAB&rs=AOn4CLBBoM3adPt3XGxO3fGKNBTLhnPD_A',
		provider: 'YOUTUBE',
		providerValue: 'I0xRbWqIohQ',
		malId: 5680,
	},
	{
		title: "K-On | Don't say lazy",
		coverUrl:
			'https://i.ytimg.com/vi/Wz-pNcgYo0c/hqdefault.jpg?sqp=-oaymwEnCOADEI4CSFryq4qpAxkIARUAAIhCGAHYAQHiAQoIGBACGAY4AUAB&rs=AOn4CLD2jCd82F5EJ-xzN5ujHruh-UvuvQ',
		provider: 'YOUTUBE',
		providerValue: 'Wz-pNcgYo0c',
		malId: 5680,
	},
	{
		title: 'K-On | Tenshi ni Fureta yo!',
		coverUrl:
			'https://i.ytimg.com/vi/y9NS5IHLunw/hq720.jpg?sqp=-oaymwEnCNAFEJQDSFryq4qpAxkIARUAAIhCGAHYAQHiAQoIGBACGAY4AUAB&rs=AOn4CLDpJO65TmF7VDfY1gLtFGkXBsy4zg',
		provider: 'YOUTUBE',
		providerValue: 'y9NS5IHLunw',
		malId: 5680,
	},
	{
		title: 'K-On | Watashi no Koi wa Hotchkiss',
		coverUrl:
			'https://i.ytimg.com/vi/4KziY05zHeQ/hq720.jpg?sqp=-oaymwEnCNAFEJQDSFryq4qpAxkIARUAAIhCGAHYAQHiAQoIGBACGAY4AUAB&rs=AOn4CLDULJkcnf9x09rsCXrcNGVYNTCXzg',
		provider: 'YOUTUBE',
		providerValue: '4KziY05zHeQ',
		malId: 5680,
	},
	{
		title: 'K-On | Listen!',
		coverUrl:
			'https://i.ytimg.com/vi/zK3YHKYM6PQ/hq720.jpg?sqp=-oaymwEnCNAFEJQDSFryq4qpAxkIARUAAIhCGAHYAQHiAQoIGBACGAY4AUAB&rs=AOn4CLDGWixZ6vCbGjovr4akEVrSJwzGmQ',
		provider: 'YOUTUBE',
		providerValue: 'zK3YHKYM6PQ',
		malId: 5680,
	},
];

export default function App() {
	const columns = process.stdout.columns;
	const rows = process.stdout.rows;

	const [selected, setSelected] = useState(0);
	const [playing, setPlaying] = useState(false);
	const [playerProcess, setPlayerProcess] =
		useState<ChildProcessByStdio<Stream.Writable, null, null>>();
	const [currentTitle, setCurrentTitle] = useState('');
	const [error, setError] = useState<string | null>(null);

	const visibleRows = rows - 12; // Adjust as needed to fit layout
	const start = Math.max(0, selected - Math.floor(visibleRows / 2));
	const end = Math.min(items.length, start + visibleRows);
	const visibleItems = items.slice(start, end);

	useInput((input, key) => {
		if (key.upArrow) {
			setSelected(prev => (prev === 0 ? items.length - 1 : prev - 1));
		} else if (key.downArrow) {
			setSelected(prev => (prev === items.length - 1 ? 0 : prev + 1));
		} else if (key.return) {
			playMusic(items[selected]);
		} else if (input === 'q') {
			playerProcess?.kill();
			process.exit(0);
		}
	});

	const playMusic = async (song: any) => {
		setPlaying(true);
		setCurrentTitle(song.title);
		const streamYtdl = await crawler.getYoutube(song.providerValue);
		const child = spawn('ffplay', ['-i', 'pipe:0', '-nodisp', '-autoexit'], {
			stdio: ['pipe', 'inherit', 'inherit'],
		});
		streamYtdl.pipe(child.stdin);

		setPlayerProcess(child);

		child.on('exit', () => setPlaying(false));
		child.on('message', msg => {
			console.log(msg);
		});
		child.on('error', err => {
			console.log(err.message);
			setError(err.message);
			setPlaying(false);
		});
	};

	return (
		<Box flexDirection="column" padding={1} width={columns} height={rows}>
			{/* Header */}
			<Box
				height={'20%'}
				justifyContent="center"
				borderStyle="round"
				borderColor="green"
				width={'100%'}
				padding={1}
			>
				<Text color="greenBright">
					🎵 Now Playing: {playing ? currentTitle + ' ▶️' : 'Nothing'}
				</Text>
			</Box>

			{/* Body */}
			<Box height={'80%'}>
				{/* Left Menu */}
				<Box
					width={20}
					flexDirection="column"
					justifyContent="space-between"
					borderStyle="round"
					borderColor="cyan"
				>
					<Box flexDirection="column">
						<Text color="cyanBright">Menu</Text>
						<Text>1. Home</Text>
						<Text>2. Search</Text>
						<Text>3. Settings</Text>
					</Box>
					<Text color="gray">Press Q to Quit</Text>
				</Box>

				{/* Right Scrollable Music List */}
				<Box
					flexGrow={1}
					flexDirection="column"
					borderStyle="round"
					borderColor="yellow"
					paddingX={1}
				>
					<Text color="yellowBright">🎶 Music List</Text>
					<Box flexDirection="column" marginTop={1}>
						{visibleItems.map((song, i) => {
							const actualIndex = start + i;
							const isSelected = actualIndex === selected;
							return (
								<Text
									key={actualIndex}
									color={isSelected ? 'greenBright' : undefined}
								>
									{isSelected ? '👉 ' : '   '}
									{song.title}
								</Text>
							);
						})}
					</Box>
				</Box>
			</Box>

			{error && (
				<Box marginTop={1}>
					<Text color="red">❌ Error: {error}</Text>
				</Box>
			)}
		</Box>
	);
}
