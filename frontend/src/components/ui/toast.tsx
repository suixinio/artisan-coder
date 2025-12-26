// Wrapper for sonner toast
import { toast as sonnerToast } from 'sonner'

type ToastProps = Parameters<typeof sonnerToast>[0]

export function toast(props: ToastProps) {
  return sonnerToast(props)
}

export { sonnerToast as toastFn, type ToastProps }
