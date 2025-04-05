import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/playground/interpolation')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/playground/interpolation"!</div>
}
