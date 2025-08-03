const formatTime = (seconds) => {
	const min = Math.floor(seconds / 60);
	const sec = Math.floor(seconds % 60);
	return `${min}:${sec.toString().padStart(2, '0')}`;
};

function stripAnsi(str) {
	return str.replace(/\x1b\[[0-9;]*m/g, '');
}

function toSlug(text) {
	return text
		.toLowerCase()
		.normalize("NFD") // remove accents
		.replace(/[\u0300-\u036f]/g, "") // remove diacritics
		.replace(/[^a-z0-9\s-]/g, "") // remove non-alphanumeric
		.trim()
		.replace(/\s+/g, "-"); // replace spaces with hyphens
}


export default {
  formatTime,
	stripAnsi,
	toSlug
}