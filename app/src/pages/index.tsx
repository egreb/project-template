import { redirect, useLoaderData } from 'react-router-dom'
import { httpClient } from '../lib/http'
import { LoaderData } from '../types/router'

interface Me {
    id: string
    username: string
}

export async function Loader() {
    const [err, res] = await httpClient.get<Me>('/api/auth/me')
    if (err) {
        return redirect('/login?message' + err)
    }

    return {
        me: res?.data,
    }
}
export default function IndexPage() {
    const loader = useLoaderData() as LoaderData<typeof Loader>
    return <div class={'bg-red-400'}>Welcome, {loader.me?.username}</div>
}
