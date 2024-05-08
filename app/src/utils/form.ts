import { ActionFunctionArgs } from 'react-router-dom'
import * as v from 'valibot'

export function form<T = unknown>(schema: v.BaseSchema<T>) {
    const parse = (data: unknown) => {
        return v.safeParse(schema, data)
    }

    const values = (data: unknown) => {
        const res = parse(data)
        if (!res.success) {
            return null
        }

        return res.output as v.Input<typeof schema>
    }

    const validate = (data: FormData) => {
        const result = parse(Object.fromEntries(data))
        if (result.success) {
            return [null, result.output as T]
        }

        const errors: Record<string, string> = {}
        for (const error of result.issues) {
            if (!error.path) {
                continue
            }

            for (const p of error.path) {
                if (typeof p.key !== 'string') {
                    continue
                }
                errors[p.key as string] = error.message
                break
            }
        }

        return [errors, result.output as T]
    }

    const request = {
        validate: async (ctx: ActionFunctionArgs) => validate(await ctx.request.formData()),
    }

    return {
        parse,
        values,
        request,
        validate,
    }
}
