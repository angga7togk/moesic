export async function getMoesicData(url) {
	const res = await fetch(url);
	const reader = res.body.getReader();
	const decoder = new TextDecoder("utf-8");

	let buffer = "";
	const result = [];

	let current = null;

	while (true) {
		const { done, value } = await reader.read();
		if (done) break;

		buffer += decoder.decode(value, { stream: true });
		const lines = buffer.split("\n");

		// Simpan sisa baris terakhir kalau belum lengkap
		buffer = lines.pop();

		for (let line of lines) {
			line = line.trim();

			if (line.startsWith("###")) {
				const title = line.replace(/^###\s*/, "").trim();
				current = { title, songs: [] };
				result.push(current);
			}

			if (line.startsWith("- [") && current) {
				const match = line.match(/\[([^\]]+)\]\(([^)]+)\)/);
				if (match) {
					const [, title, url] = match;
					current.songs.push({ title, url });
				}
			}
		}
	}

	// Proses sisa buffer terakhir (jika ada)
	if (buffer) {
		const line = buffer.trim();

		if (line.startsWith("###")) {
			const title = line.replace(/^###\s*/, "").trim();
			current = { title, songs: [] };
			result.push(current);
		}

		if (line.startsWith("- [") && current) {
			const match = line.match(/\[([^\]]+)\]\(([^)]+)\)/);
			if (match) {
				const [, title, url] = match;
				current.songs.push({ title, url });
			}
		}
	}
	return result;
}


export default {
  getMoesicData
}
