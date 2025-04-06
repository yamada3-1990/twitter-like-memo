const SERVER_URL = import.meta.env.VITE_BACKEND_URL || 'http://127.0.0.1:9000';

export function Memo(id, title, body, tags) {
    this.id = id;
    this.title = title;
    this.body = body;
    this.tags = tags;
};

export const MemoListResponse = {
    memos: Memo
};

export const fetchMemos = async () => {
    const response = await fetch(`${SERVER_URL}/memos`, {
        method: 'GET',
        mode: 'cors',
        headers: {
            'Content-Type': 'application/json',
            Accept: 'application/json',
        },
    });

    if (response.status >= 400) {
        const error = await response.json();
        throw new Error(`Failed to fetch memos from the server: ${JSON.stringify(error)}`);
    }
    return response.json();
};

export const addMemos = async (input) => {
    const data = new FormData();
    data.append('title', input.title || '無題');
    data.append('body', input.body);
    data.append('tags', input.tags || '');

    // FormDataの内容を確認
    console.log('FormData contents:');
    for (let [key, value] of data.entries()) {
        console.log(`${key}: ${value}`);
    }

    const response = await fetch(`${SERVER_URL}/memos`, {
        method: 'POST',
        mode: 'cors',
        headers: {
            'Accept': 'application/json',
        },
        body: data,
    });

    if (response.status >= 400) {
        const contentType = response.headers.get('content-type');
        let errorMessage;
        
        if (contentType && contentType.includes('application/json')) {
            const errorJson = await response.json();
            errorMessage = JSON.stringify(errorJson);
        } else {
            errorMessage = await response.text();
        }
        
        throw new Error(`Failed to post memos to the server: ${errorMessage}`);
    }

    return response.json();
};

export const deleteMemo = async (memoId) => {
    const response = await fetch(`${SERVER_URL}/memos/${memoId}`, {
        method: 'DELETE',
        mode: 'cors',
        headers: {
            'Content-Type': 'application/json',
            Accept: 'application/json',
        },
    });

    if (response.status >= 400) {
        const error = await response.json();
        throw new Error(`Failed to delete memo: ${JSON.stringify(error)}`);
    }
    return response.json();
}

export const searchMemosByKeyword = async (keyword) => {
    const queryParams = new URLSearchParams();
    if (keyword.keyword) {
        queryParams.append('keyword', keyword.keyword);
    }

    const response = await fetch(`${SERVER_URL}/search/keyword?${queryParams.toString()}`, {
        method: 'GET',
        mode: 'cors',
        headers: {
            'Content-Type': 'application/json',
            Accept: 'application/json',
        },
    });

    if (response.status >= 400) {
        const error = await response.json();
        throw new Error(`Failed to search memos: ${JSON.stringify(error)}`);
    }
    return response.json();
}

export const searchMemosByTags = async (tags) => {
    const queryParams = new URLSearchParams();
    if (tags.tags) {
        queryParams.append('tags', tags.tags);
    }

    const response = await fetch(`${SERVER_URL}/search/tags?${queryParams.toString()}`, {
        method: 'GET',
        mode: 'cors',
        headers: {
            'Content-Type': 'application/json',
            Accept: 'application/json',
        },
    });

    if (response.status >= 400) {
        const error = await response.json();
        throw new Error(`Failed to search memos: ${JSON.stringify(error)}`);
    }
    return response.json();
}