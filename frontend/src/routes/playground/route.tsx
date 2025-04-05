import { createFileRoute, Outlet } from "@tanstack/react-router";
import { ArrowLeft } from "lucide-react";
import { Link } from "@tanstack/react-router";

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
  return (
    <div className="container mx-auto py-10">
      <div className="mb-8 flex items-center">
        <Button variant="ghost" size="sm" asChild className="mr-2">
          <Link to="/">
            <ArrowLeft className="mr-2 h-4 w-4" />
            Back to Home
          </Link>
        </Button>
        <h1 className="font-comic text-3xl font-bold">
          Numerical Methods Playground
        </h1>
      </div>

      <Card className="">
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
          <Tabs>
            <TabsList className="grid w-full grid-cols-4 border-2 border-black p-0">
              <TabsTrigger
                value="root-finding"
                className="font-comic data-[state=active]:bg-primary data-[state=active]:text-primary-foreground h-full"
                asChild
              >
                <Link to={"/playground/root-finding"}>Root Finding</Link>
              </TabsTrigger>
              <TabsTrigger
                value="integration"
                className="font-comic data-[state=active]:bg-primary data-[state=active]:text-white"
              >
                Integration
              </TabsTrigger>
              <TabsTrigger
                value="differential"
                className="font-comic data-[state=active]:bg-primary data-[state=active]:text-white"
              >
                Differential Eqs
              </TabsTrigger>
              <TabsTrigger
                value="interpolation"
                className="font-comic data-[state=active]:bg-primary data-[state=active]:text-white"
              >
                Interpolation
              </TabsTrigger>
            </TabsList>
            <TabsContent value="root-finding">
              <Outlet />
            </TabsContent>
            <TabsContent value="integration">
              <div>integration</div>
            </TabsContent>
            <TabsContent value="differential">
              <div>differential</div>
            </TabsContent>
            <TabsContent value="interpolation">
              <div>interpolation</div>
            </TabsContent>
          </Tabs>
        </CardContent>
      </Card>
    </div>
  );
}
