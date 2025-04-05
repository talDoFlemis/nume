import { createFileRoute } from '@tanstack/react-router'

export const Route:unknown = createFileRoute('/playground')({
  component: About,
})

function About() {
  return <div className="p-2">Hello from About!</div>
}