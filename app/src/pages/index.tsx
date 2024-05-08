import { ActionFunctionArgs, Form, redirect, useLoaderData } from 'react-router-dom'
import { httpClient } from '../lib/http'
import { LoaderData } from '../types/router'
import * as v from 'valibot'
import { form } from '../utils/form'

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

const formSchema = v.variant('_intent', [
    v.object({
        _intent: v.literal('signout'),
    }),
    v.object({
        _intent: v.literal('empty'),
    }),
])

export async function Action(ctx: ActionFunctionArgs) {
    const data = Object.fromEntries(await ctx.request.formData())
    if (!data) {
        return null
    }

    if (v.is(formSchema, data)) {
        if (data._intent === 'signout') {
            await httpClient.post('/api/auth/signout', null)

            return redirect('/login')
        }
    }
}
export default function IndexPage() {
    const loader = useLoaderData() as LoaderData<typeof Loader>
    return (
        <div class={'h-dvh bg-slate-50 text-gray-800 flex flex-col gap-y-4'}>
            <header>Welcome, {loader.me?.username}</header>

            <Form method="POST">
                <button
                    type="submit"
                    name="_intent"
                    value="signout"
                    class={'bg-blue-400 hover:bg-blue-500 px-4 py-2 rounded-md text-center text-white'}
                >
                    Signout
                </button>
            </Form>
        </div>
    )
}
