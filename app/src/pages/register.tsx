import { ActionFunctionArgs, Form } from 'react-router-dom'
import { form } from '../utils/form'
import * as v from 'valibot'
import { httpClient } from '../lib/http'
import { Path, redirect } from '../router'

const schema = v.object({
    username: v.string(),
    password: v.string(),
})

type RegisterSchema = v.Input<typeof schema>

export async function Action(ctx: ActionFunctionArgs) {
    const formData = await ctx.request.formData()
    if (!formData) {
        return null
    }

    // TODO: Validate this
    const data = form(schema).parse(Object.fromEntries(formData))

    if (!data.success) {
        return {
            error: true,
        }
    }

    const [error, res] = await httpClient.post<RegisterSchema, string>('/api/auth/register', data.output)
    if (error) {
        return {
            error: true,
        }
    }

    return redirect(('/login?' + res?.data) as Path)
}

export default function RegisterPage() {
    return (
        <div class={'bg-green-400 h-dvh grid place-items-center'}>
            <div class={'max-w-md mx-auto rounded-md p-8 bg-white'}>
                <h1 class={'text-xl'}>Register</h1>
                <Form method="POST" class="pt-8 flex flex-col gap-y-3 rounded-md bg-white">
                    <label for="username" class={'flex flex-col gap-y-3'}>
                        Username
                        <input type="text" id="username" name="username" />
                    </label>
                    <label for="password" class={'flex flex-col gap-y-3'}>
                        Password
                        <input type="password" id="password" name="password" />
                    </label>

                    <button type="submit" class={'px-4 py-2 bg-slate-800 text-slate-100 mt-3'}>
                        Register
                    </button>
                </Form>
            </div>
        </div>
    )
}
