/**
 * View Transitions API 的类型声明。
 *
 * 该 API 目前还没进入 TypeScript 的内置 DOM lib，这里补一份最小定义。
 * startViewTransition 声明为可选，调用前必须做能力检测。
 */
export {}

declare global {
  interface ViewTransition {
    readonly ready: Promise<void>
    readonly finished: Promise<void>
    readonly updateCallbackDone: Promise<void>
    skipTransition(): void
  }

  interface Document {
    startViewTransition?: (callback: () => void | Promise<void>) => ViewTransition
  }
}
