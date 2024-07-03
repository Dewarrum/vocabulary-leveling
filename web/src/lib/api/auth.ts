import { createQuery } from "@tanstack/svelte-query";

export type Profile = {
    name: string;
    username: string;
    roles: string[];
}

export async function getProfile() {
    const response = await fetch('/auth/profile', {
        method: 'GET'
    });

    if (response.status === 401) {
        return null;
    }

    const data: Profile = await response.json();

    return data;
}

export const createProfileQuery = () => createQuery({
    queryKey: ['profile'],
    queryFn: () => getProfile()
});