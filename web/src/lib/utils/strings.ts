export function isPartOfURL(urlA: string, urlB: string, exact = false): boolean {
	if (exact) {
		return urlA == urlB;
	}

	return urlA.startsWith(urlB);
}
