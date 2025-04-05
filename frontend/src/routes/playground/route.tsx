import {
  createFileRoute,
  Outlet,
  useLocation,
  useNavigate,
} from "@tanstack/react-router";
import { ArrowLeft } from "lucide-react";
import { Link } from "@tanstack/react-router";
import React from "react";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";

export const Route = createFileRoute("/playground")({
  component: NumericalMethodsPlayground,
});

function NumericalMethodsPlayground() {
  const location = useLocation();
  const navigate = useNavigate({ from: "/playground" });

  const currentTab = location.pathname.split("/").pop() || "root-finding";

  React.useEffect(() => {
    if (location.pathname === "/playground") {
      navigate({
        to: "/playground/root-finding",
        search: {
          method: "bisection",
          delta: 1,
          initialGuess: 0,
          iterations: 10,
          fn: "x^2 - 4",
          error: 0.01,
        },
      });
    }
  }, [location.pathname, navigate]);

  return (
    <div className="container mx-auto py-10">
      <div className="mb-8 flex items-center">
        <Button variant="ghost" size="sm" asChild className="mr-2">
          <Link to="/">
            <ArrowLeft className="h -4 mr-2 w-4" />
            Back to Home
          </Link>
        </Button>
        <h1 className="font-comic text-3xl font-bold">
          Numerical Methods Playground
        </h1>
      </div>

      <Card>
        <div className="absolute -top-20 -right-20 h-40 w-40 rounded-full bg-yellow-300 opacity-10"></div>
        <CardHeader>
          <CardTitle className="font-comic text-2xl">
            Interactive Playground
          </CardTitle>
          <CardDescription>
            Experiment with different numerical methods and see how they work in
            real-time
          </CardDescription>
        </CardHeader>
        <CardContent>
          <Tabs value={currentTab}>
            <TabsList className="flex h-auto w-full flex-wrap items-center justify-start border-2 border-black p-0 md:space-y-0">
              <TabsTrigger
                value="root-finding"
                className="font-comic data-[state=active]:bg-primary data-[state=active]:text-primary-foreground h-full"
                asChild
              >
                <Link
                  to={"/playground/root-finding"}
                  search={{
                    method: "bisection",
                    delta: 1,
                    initialGuess: 0,
                    iterations: 10,
                    fn: "x^2 - 4",
                    error: 0.01,
                  }}
                >
                  Root Finding
                </Link>
              </TabsTrigger>
              <TabsTrigger
                value="integration"
                className="font-comic data-[state=active]:bg-primary data-[state=active]:text-white"
                asChild
              >
                <Link to="/playground/integration">Integration</Link>
              </TabsTrigger>
              <TabsTrigger
                value="differential"
                className="font-comic data-[state=active]:bg-primary data-[state=active]:text-white"
                asChild
              >
                <Link to="/playground/differential">Differential Eqs</Link>
              </TabsTrigger>
              <TabsTrigger
                value="interpolation"
                className="font-comic data-[state=active]:bg-primary data-[state=active]:text-white"
                asChild
              >
                <Link to="/playground/interpolation">Interpolation</Link>
              </TabsTrigger>
            </TabsList>
            <TabsContent value="root-finding">
              <Outlet />
            </TabsContent>
            <TabsContent value="integration">
              <Outlet />
            </TabsContent>
            <TabsContent value="differential">
              <Outlet />
            </TabsContent>
            <TabsContent value="interpolation">
              <Outlet />
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>
    </div>
  );
}
