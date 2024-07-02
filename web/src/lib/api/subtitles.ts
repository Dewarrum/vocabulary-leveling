export type Subtitle = {
    id: string;
    videoName: string;
    startMs: number;
    endMs: number;
    text: string;
}

async function searchSubtitles(query: string) {
    if (!query) {
        return [];
    }

    const response = await fetch(`/api/subtitles/search?query=${query}`, {
        method: 'GET'
    });

    const data: Subtitle[] = await response.json();

    return data;
}

export { searchSubtitles };