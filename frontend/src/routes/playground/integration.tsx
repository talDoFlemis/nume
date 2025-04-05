import { createFileRoute } from "@tanstack/react-router";

export const Route = createFileRoute("/playground/integration")({
  component: RouteComponent,
});

function RouteComponent() {
  return <div>Hello "/playground/integration"!</div>;
}
