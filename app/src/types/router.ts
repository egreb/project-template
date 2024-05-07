import { type ActionFunction, type LoaderFunction } from 'react-router-dom'

export type LoaderError = Record<string, string | undefined>
export type LoaderData<TLoaderFn extends LoaderFunction> =
    Awaited<ReturnType<TLoaderFn>> extends Response | infer D ? D : never
export type ActionError = Record<string, string | undefined>
export type ActionData<TActionFn extends ActionFunction> =
    Awaited<ReturnType<TActionFn>> extends Response | infer D ? D : never
