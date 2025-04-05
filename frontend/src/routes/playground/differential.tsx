import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/playground/differential')({
  component: RouteComponent,
})

function RouteComponent() {
  return <div>Hello "/playground/differential"!</div>
}
