import { PUBLIC_API } from "$env/static/public";

export type Subtitle = {
    id: string;
    videoId: string;
    sequence: number;
    startMs: number;
    endMs: number;
    text: string;
}

async function searchSubtitles(query: string) {
    if (!query) {
        return [];
    }

    const response = await fetch(`${PUBLIC_API}/api/subtitles/search?query=${query}`, {
        method: 'GET'
    });
    console.log(response);

    const data: Subtitle[] = await response.json();
    console.log(data);

    return data;
}

export { searchSubtitles };