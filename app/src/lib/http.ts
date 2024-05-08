const baseConfig: Partial<RequestInit> = {
    credentials: 'include',
    headers: {
        'Content-Type': 'application/json',
    },
}

type Options = Omit<RequestInit, 'method'>
type URL = `/${string}`
interface HttpClientResponse<T = Record<string, unknown>> extends Response {
    data: T | null
}

type ResponseError<T> = {
    status: number
    statusText: string
    data: T | null
}

async function getJSON<T>(res: Response): Promise<T | null> {
    try {
        return (await res.json()) as T
    } catch {
        return null
    }
}

class HttpClient {
    constructor() {}

    public async get<T = unknown, E = ResponseError<unknown>>(
        url: URL,
        opts?: Options
    ): Promise<[null, HttpClientResponse<T>] | [E, null] | [E, HttpClientResponse<T>]> {
        try {
            const res = await fetch(url, {
                ...opts,
                ...baseConfig,
                method: 'GET',
            })

            const data = await getJSON<T>(res)
            if (!res.ok) {
                const { status, statusText, headers } = res
                return [
                    { status, statusText, data } as E,
                    { status, statusText, headers, data } as HttpClientResponse<T>,
                ]
            }

            const response: Partial<HttpClientResponse<T>> = res
            response.data = data

            return [null, response as HttpClientResponse<T>]
        } catch (e) {
            console.log({ e })
            return [e as E, null]
        }
    }

    public async post<T = unknown, R = unknown, E = unknown>(
        url: URL,
        data: T,
        opts?: Options
    ): Promise<[null, HttpClientResponse<R>] | [E, null]> {
        try {
            const body = JSON.stringify(data)
            const res = await fetch(url, {
                ...opts,
                ...baseConfig,
                method: 'POST',
                body: body,
            })

            const resData = await getJSON<R>(res)
            if (!res.ok) {
                return [resData as E, null]
            }

            const response: Partial<HttpClientResponse<R>> = res
            response.data = resData

            return [null, response as HttpClientResponse<R>]
        } catch (e) {
            return [e as E, null]
        }
    }
}

export const httpClient = new HttpClient()
