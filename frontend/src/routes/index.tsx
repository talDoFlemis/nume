import { ArrowRight, Code, Calculator, LineChart, Zap } from "lucide-react";
import { siGithub } from "simple-icons";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { ModeToggle } from "@/components/theme/mode-toggle-btn";

import { createFileRoute, Link } from "@tanstack/react-router";

export const Route: unknown = createFileRoute("/")({
  component: Index,
});

function Index() {
  return (
    <div className="bg-background flex min-h-screen flex-col scroll-smooth">
      <header className="bg-background/95 supports-[backdrop-filter]:bg-background/60 sticky top-0 z-40 w-full border-b backdrop-blur">
        <div className="container mx-auto flex h-16 items-center space-x-4 sm:justify-between sm:space-x-0">
          <div className="flex gap-6 md:gap-10">
            <a href="/" className="flex items-center space-x-2">
              <Calculator className="h-6 w-6" />
              <span className="inline-block font-bold">nume</span>
            </a>
            <nav className="hidden gap-6 md:flex">
              <a
                href="#features"
                className="hover:text-primary flex items-center text-lg font-medium transition-colors"
              >
                Features
              </a>
              <a
                href="#methods"
                className="hover:text-primary flex items-center text-lg font-medium transition-colors"
              >
                Methods
              </a>
              <a
                href="#playground"
                className="hover:text-primary flex items-center text-lg font-medium transition-colors"
              >
                Playground
              </a>
            </nav>
          </div>
          <div className="flex flex-1 items-center justify-end space-x-4">
            <nav className="flex items-center space-x-2">
              <ModeToggle />
              <a
                href="https://github.com/taldoflemis/nume"
                target="_blank"
                rel="noreferrer"
              >
                <div className="bg-background hover:bg-muted hover:text-primary focus-visible:ring-ring inline-flex h-9 w-9 items-center justify-center rounded-md text-sm font-medium transition-colors focus-visible:ring-1 focus-visible:outline-none disabled:pointer-events-none disabled:opacity-50">
                  <svg
                    role="img"
                    viewBox="0 0 24 24"
                    className="h-5 w-5"
                    fill="currentColor"
                  >
                    <path d={siGithub.path} />
                  </svg>
                  <span className="sr-only">GitHub</span>
                  <span className="sr-only">GitHub</span>
                </div>
              </a>
              <Button asChild>
                <a href="#playground">Try It Now</a>
              </Button>
            </nav>
          </div>
        </div>
      </header>
      <main className="flex-1">
        <section className="space-y-6 pt-6 pb-8 md:pt-10 md:pb-12 lg:py-32">
          <div className="container mx-auto flex max-w-[64rem] flex-col items-center gap-4 text-center">
            <a
              href="https://github.com/taldoflemis/nume"
              className="bg-muted rounded-2xl px-4 py-1.5 text-sm font-medium"
              target="_blank"
            >
              Follow along on GitHub
            </a>
            <h1 className="font-heading text-3xl sm:text-5xl md:text-6xl lg:text-7xl">
              Numerical Methods Made{" "}
              <span className="text-primary">Simple</span>
            </h1>
            <p className="text-muted-foreground max-w-[42rem] leading-normal sm:text-xl sm:leading-8">
              Explore, learn, and visualize numerical methods with our
              interactive playground. Perfect for students, educators, and
              professionals in computational mathematics.
            </p>
            <div className="space-x-4">
              <Button asChild size="lg">
                <a href="#playground">
                  Get Started
                  <ArrowRight className="ml-2 h-4 w-4" />
                </a>
              </Button>
              <Button variant="outline" size="lg" asChild>
                <a href="#methods">Learn More</a>
              </Button>
            </div>
          </div>
        </section>
        <section
          id="features"
          className="container mx-auto space-y-6 py-8 md:py-12 lg:py-24"
        >
          <div className="mx-auto flex max-w-[58rem] flex-col items-center space-y-4 text-center">
            <h2 className="font-heading text-3xl leading-[1.1] sm:text-3xl md:text-6xl">
              Features
            </h2>
            <p className="text-muted-foreground max-w-[85%] leading-normal sm:text-lg sm:leading-7">
              nume provides a comprehensive platform for understanding and
              applying numerical methods
            </p>
          </div>
          <div className="mx-auto grid grid-cols-1 justify-center gap-4 px-6 md:max-w-[64rem] md:grid-cols-3">
            <Card className="hover:translate-x-[3px] hover:translate-y-[3px] hover:shadow-[5px_5px_0px_0px]">
              <div className="bg-primary absolute -top-20 -right-20 h-40 w-40 rounded-full opacity-20"></div>
              <CardHeader className="p-6">
                <CardTitle>Interactive Visualizations</CardTitle>
                <CardDescription>
                  See numerical methods in action with dynamic, interactive
                  visualizations that help you understand the underlying
                  concepts.
                </CardDescription>
              </CardHeader>
              <CardFooter className="flex justify-end">
                <LineChart className="text-primary" />
              </CardFooter>
            </Card>
            <Card className="hover:translate-x-[3px] hover:translate-y-[3px] hover:shadow-[5px_5px_0px_0px]">
              <div className="bg-secondary/20 absolute -top-20 -right-20 h-40 w-40 rounded-full"></div>
              <CardHeader className="p-6">
                <CardTitle>Step-by-Step Solutions</CardTitle>
                <CardDescription>
                  Follow along with detailed step-by-step solutions that break
                  down complex numerical processes.
                </CardDescription>
              </CardHeader>
              <CardFooter className="flex justify-end">
                <Code className="text-secondary h-6 w-6" />
              </CardFooter>
            </Card>
            <Card className="hover:translate-x-[3px] hover:translate-y-[3px] hover:shadow-[5px_5px_0px_0px]">
              <div className="bg-accent/20 absolute -top-20 -right-20 h-40 w-40 rounded-full"></div>
              <CardHeader className="p-6">
                <CardTitle>Real-time Computation</CardTitle>
                <CardDescription>
                  Experiment with different parameters and see results instantly
                  with our high-performance computation engine.
                </CardDescription>
              </CardHeader>
              <CardFooter className="flex justify-end">
                <Zap className="text-primary h-6 w-6" />
              </CardFooter>
            </Card>
          </div>
        </section>
        <section
          id="methods"
          className="container mx-auto space-y-6 py-8 md:py-12 lg:py-24"
        >
          <div className="container mx-auto flex flex-col items-center space-y-4 text-center">
            <h2 className="font-heading text-3xl leading-[1.1] sm:text-3xl md:text-6xl">
              Numerical Methods
            </h2>
            <p className="text-muted-foreground max-w-[85%] leading-normal sm:text-lg sm:leading-7">
              Explore a variety of numerical methods for solving mathematical
              problems
            </p>
          </div>
          <div className="mx-auto grid grid-cols-1 justify-center gap-12 px-6 md:max-w-[64rem] md:grid-cols-3">
            <Card className="hover:translate-x-[3px] hover:translate-y-[3px] hover:shadow-[5px_5px_0px_0px]">
              <CardHeader>
                <CardTitle>Root Finding</CardTitle>
                <CardDescription>
                  Bisection Method, Newton-Raphson Method, Secant Method, and
                  more for finding roots of equations.
                </CardDescription>
              </CardHeader>
              <CardFooter className="flex items-center justify-end">
                <Button variant="outline" size="sm" asChild>
                  <a href="/methods/root-finding">
                    Explore
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </a>
                </Button>
              </CardFooter>
            </Card>
            <Card className="hover:translate-x-[3px] hover:translate-y-[3px] hover:shadow-[5px_5px_0px_0px]">
              <CardHeader>
                <CardTitle>Numerical Integration</CardTitle>
                <CardDescription>
                  Trapezoidal Rule, Simpson's Rule, and Gaussian Quadrature for
                  approximating definite integrals.
                </CardDescription>
              </CardHeader>
              <CardFooter className="flex items-center justify-end">
                <Button variant="outline" size="sm" asChild>
                  <a href="/methods/integration">
                    Explore
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </a>
                </Button>
              </CardFooter>
            </Card>
            <Card className="hover:translate-x-[3px] hover:translate-y-[3px] hover:shadow-[5px_5px_0px_0px]">
              <CardHeader>
                <CardTitle>Differential Equations</CardTitle>
                <CardDescription>
                  Euler's Method, Runge-Kutta Methods, and more for solving
                  ordinary differential equations.
                </CardDescription>
              </CardHeader>
              <CardFooter className="flex items-center justify-end">
                <Button variant="outline" size="sm" asChild>
                  <a href="/methods/differential-equations">
                    Explore
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </a>
                </Button>
              </CardFooter>
            </Card>
            <Card className="hover:translate-x-[3px] hover:translate-y-[3px] hover:shadow-[5px_5px_0px_0px]">
              <CardHeader>
                <CardTitle>Linear Systems</CardTitle>
                <CardDescription>
                  Gaussian Elimination, LU Decomposition, and Iterative Methods
                  for solving systems of linear equations.
                </CardDescription>
              </CardHeader>
              <CardFooter className="flex items-center justify-end">
                <Button variant="outline" size="sm" asChild>
                  <a href="/methods/linear-systems">
                    Explore
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </a>
                </Button>
              </CardFooter>
            </Card>
            <Card className="hover:translate-x-[3px] hover:translate-y-[3px] hover:shadow-[5px_5px_0px_0px]">
              <CardHeader>
                <CardTitle>Interpolation</CardTitle>
                <CardDescription>
                  Lagrange Polynomials, Newton's Divided Differences, and Spline
                  Interpolation for fitting curves to data points.
                </CardDescription>
              </CardHeader>
              <CardFooter className="flex items-center justify-end">
                <Button variant="outline" size="sm" asChild>
                  <a href="/methods/interpolation">
                    Explore
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </a>
                </Button>
              </CardFooter>
            </Card>
            <Card className="hover:translate-x-[3px] hover:translate-y-[3px] hover:shadow-[5px_5px_0px_0px]">
              <CardHeader>
                <CardTitle>Optimization</CardTitle>
                <CardDescription>
                  Golden Section Search, Gradient Descent, and more for finding
                  minima and maxima of functions.
                </CardDescription>
              </CardHeader>
              <CardFooter className="flex items-center justify-end">
                <Button variant="outline" size="sm" asChild>
                  <a href="/methods/optimization">
                    Explore
                    <ArrowRight className="ml-2 h-4 w-4" />
                  </a>
                </Button>
              </CardFooter>
            </Card>
          </div>
        </section>
        <section
          id="playground"
          className="container mx-auto space-y-6 py-8 md:py-12 lg:py-24"
        >
          <div className="mx-auto flex max-w-[58rem] flex-col items-center space-y-4 text-center">
            <h2 className="font-heading text-3xl leading-[1.1] sm:text-3xl md:text-6xl">
              Interactive Playground
            </h2>
            <p className="text-muted-foreground max-w-[85%] leading-normal sm:text-lg sm:leading-7">
              Try out numerical methods with our interactive playground
            </p>
          </div>
          <div className="mx-auto flex justify-center">
            <div className="bg-background relative w-full max-w-4xl overflow-hidden rounded-lg border p-4">
              <div className="bg-muted flex aspect-video items-center justify-center rounded-md">
                <div className="space-y-4 text-center">
                  <h3 className="text-2xl font-bold">Playground Preview</h3>
                  <p className="text-muted-foreground">
                    Interactive visualization and computation of numerical
                    methods
                  </p>
                  <Button asChild variant="cartoon">
                    <Link to="/playground">
                      Launch Playground
                      <ArrowRight className="ml-2 h-4 w-4" />
                    </Link>
                  </Button>
                </div>
              </div>
            </div>
          </div>
        </section>
        <section className="container mx-auto py-8 md:py-12 lg:py-24">
          <div className="mx-auto max-w-[58rem] space-y-12">
            <div className="flex flex-col items-center justify-center gap-4 text-center">
              <h2 className="font-heading text-3xl leading-[1.1] sm:text-3xl md:text-6xl">
                Ready to dive in?
              </h2>
              <p className="text-muted-foreground max-w-[85%] leading-normal sm:text-lg sm:leading-7">
                Start exploring numerical methods with our interactive tools
              </p>
              <Button size="lg" variant="cartoon" asChild>
                <a href="/playground">
                  Get Started
                  <ArrowRight className="ml-2 h-4 w-4" />
                </a>
              </Button>
            </div>
          </div>
        </section>
      </main>
      <footer className="border-t py-6 md:py-0">
        <div className="container mx-auto flex flex-col items-center justify-between gap-4 md:h-24 md:flex-row">
          <div className="flex flex-col items-center gap-4 px-8 md:flex-row md:gap-2 md:px-0">
            <Calculator className="h-6 w-6" />
            <p className="text-center text-sm leading-loose md:text-left">
              &copy; {new Date().getFullYear()} nume. All rights reserved.
            </p>
          </div>
          <div className="flex items-center space-x-4">
            <a
              href="https://github.com/taldoflemis/nume"
              target="_blank"
              rel="noreferrer"
            >
              <div className="bg-background hover:bg-muted hover:text-primary focus-visible:ring-ring inline-flex h-9 w-9 items-center justify-center rounded-md text-sm font-medium transition-colors focus-visible:ring-1 focus-visible:outline-none disabled:pointer-events-none disabled:opacity-50">
                <svg
                  role="img"
                  viewBox="0 0 24 24"
                  className="h-5 w-5"
                  fill="currentColor"
                >
                  <path d={siGithub.path} />
                </svg>
                <span className="sr-only">GitHub</span>
              </div>
            </a>
          </div>
        </div>
      </footer>
    </div>
  );
}
