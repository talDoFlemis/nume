import { createFileRoute, useNavigate } from "@tanstack/react-router";
import { useState } from "react";
import { LineChart, Calculator } from "lucide-react";
import { zodValidator } from "@tanstack/zod-adapter";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";

import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";

type Method = "bisection" | "newthon-raphson" | "secant";

const methods: Method[] = ["bisection", "newthon-raphson", "secant"];

const rootFindingSearchSchema = z.object({
  method: z.enum(["bisection", "newthon-raphson", "secant"] as const),
  fn: z.string(),
  iterations: z.coerce.number().min(1).max(100),
  delta: z.coerce.number().positive().or(z.number().negative()),
  initialGuess: z.coerce.number(),
  error: z.coerce.number().nonnegative().gt(0),
});

type RootFindingSearchSchema = z.infer<typeof rootFindingSearchSchema>;

export const Route = createFileRoute("/playground/root-finding")({
  component: RouteFinding,
  validateSearch: zodValidator(rootFindingSearchSchema),
});

function RouteFinding() {
  const search = Route.useSearch();
  const form = useForm<RootFindingSearchSchema>({
    resolver: zodResolver(rootFindingSearchSchema),
    defaultValues: {
      method: search.method,
      fn: search.fn,
      iterations: search.iterations,
      delta: search.delta,
      initialGuess: search.initialGuess,
      error: search.error,
    },
  });
  const navigate = useNavigate({ from: Route.fullPath });
  const [result, setResult] = useState<string | null>(null);
  const [isCalculating, setIsCalculating] = useState(false);

  const sleep = (ms: number) =>
    new Promise((resolve) => setTimeout(resolve, ms));

  const handleCalculate = async (data: RootFindingSearchSchema) => {
    const { method, fn, iterations, delta, initialGuess } = data;

    await sleep(1000); // Simulate a delay for the calculation
    setResult(
      `Method: ${method}\nFunction: ${fn}\nParameters: [delta: ${delta.toString()}, ${initialGuess.toString()}]\nIterations: ${iterations.toString()}\n\nResult would be displayed here with step-by-step calculations and visualization.`,
    );
  };

  const convertMethodToString = (method: Method) => {
    switch (method) {
      case "bisection":
        return "Bisection Method";
      case "newthon-raphson":
        return "Newton-Raphson Method";
      case "secant":
        return "Secant Method";
      default:
        return "Unknown Method";
    }
  };

  const onSubmit = async (data: RootFindingSearchSchema) => {
    const { method, fn, iterations, delta, initialGuess, error } = data;

    setIsCalculating(true);

    try {
      await handleCalculate(data);
    } catch (error) {
      console.error("Error during calculation:", error);
      setResult("An error occurred during calculation.");
    } finally {
      setIsCalculating(false);
    }

    await navigate({
      search: {
        method,
        fn,
        iterations,
        delta,
        initialGuess,
        error,
      },
    });
  };

  return (
    <div className="mt-6 grid gap-12 py-6 lg:grid-cols-2">
      <Card className="py-6 transition-none">
        <div className="absolute -top-20 -left-20 h-40 w-40 rounded-full bg-blue-300 opacity-10"></div>
        <CardContent>
          <Form {...form}>
            <form
              className="space-y-4"
              onSubmit={(e) => {
                e.preventDefault();
                //eslint-disable-next-line
                form.handleSubmit(onSubmit)(e);
              }}
            >
              <FormField
                control={form.control}
                name="method"
                render={({ field }) => (
                  <FormItem className="space-y-2">
                    <FormLabel className="font-comic">Method</FormLabel>
                    <Select value={field.value} onValueChange={field.onChange}>
                      <FormControl>
                        <SelectTrigger className="border-2 border-black">
                          <SelectValue placeholder="Select method" />
                        </SelectTrigger>
                      </FormControl>
                      <SelectContent className="group border-card-foreground bg-card relative overflow-hidden rounded-lg border-4 p-2 shadow-[8px_8px_0px_0px]">
                        {methods.map((method) => (
                          <SelectItem key={method} value={method}>
                            {convertMethodToString(method)}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="fn"
                render={({ field }) => (
                  <FormItem className="space-y-2">
                    <FormLabel className="font-comic">Function f(x)</FormLabel>
                    <FormControl>
                      <Input
                        id="function"
                        placeholder="e.g. x^2 - 4"
                        className="border-2 border-black"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <div className="grid grid-cols-2 gap-4">
                <FormField
                  control={form.control}
                  name="initialGuess"
                  render={({ field }) => (
                    <FormItem className="space-y-2">
                      <FormLabel className="font-comic">
                        Initial Guess
                      </FormLabel>
                      <FormControl>
                        <Input
                          id="initialGuess"
                          placeholder="e.g. 0"
                          className="border-2 border-black"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
                <FormField
                  control={form.control}
                  name="delta"
                  render={({ field }) => (
                    <FormItem className="space-y-2">
                      <FormLabel className="font-comic">Delta</FormLabel>
                      <FormControl>
                        <Input
                          id="delta"
                          placeholder="e.g. 1"
                          className="border-2 border-black"
                          {...field}
                        />
                      </FormControl>
                      <FormMessage />
                    </FormItem>
                  )}
                />
              </div>
              <div className="grid grid-cols-2 gap-4">
              <FormField
                control={form.control}
                name="iterations"
                render={({ field }) => (
                  <FormItem className="space-y-2">
                    <FormLabel className="font-comic">Iterations</FormLabel>
                    <FormControl>
                      <Input
                        id="iterations"
                        placeholder="e.g. 10"
                        className="border-2 border-black"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              <FormField
                control={form.control}
                name="error"
                render={({ field }) => (
                  <FormItem className="space-y-2">
                    <FormLabel className="font-comic">Error</FormLabel>
                    <FormControl>
                      <Input
                        placeholder="0.01"
                        className="border-2 border-black"
                        {...field}
                      />
                    </FormControl>
                    <FormMessage />
                  </FormItem>
                )}
              />
              </div>
              <Button
                variant={"cartoon"}
                className="w-full cursor-pointer"
                disabled={isCalculating}
              >
                {isCalculating ? "Calculating..." : "Calculate"}
                <Calculator className="ml-2 h-4 w-4" />
              </Button>
            </form>
          </Form>
        </CardContent>
        <div className="absolute -top-2 -left-2 h-6 w-6 rounded-full border-4 border-black bg-white"></div>
        <div className="absolute -right-2 -bottom-2 h-6 w-6 rounded-full border-4 border-black bg-white"></div>
      </Card>

      <Card>
        <div className="absolute -top-20 -right-20 h-40 w-40 rounded-full bg-green-300 opacity-10"></div>
        <CardContent className="pt-6">
          <div className="space-y-4">
            <div className="relative flex aspect-square items-center justify-center rounded-md border-2 border-black bg-gradient-to-br from-blue-100 to-purple-100">
              <div className="font-comic absolute top-2 left-2 -rotate-6 transform rounded-lg border-2 border-black bg-yellow-300 px-3 py-1 text-xs font-bold">
                Visualization
              </div>
              <LineChart className="text-primary h-16 w-16" />
              <span className="sr-only">Visualization area</span>
            </div>

            <div className="space-y-2">
              <Label htmlFor="result" className="font-comic">
                Result
              </Label>
              <Textarea
                id="result"
                value={result ?? "Results will appear here after calculation"}
                readOnly
                className="min-h-[200px] border-2 border-black font-mono text-sm"
              />
            </div>
          </div>
        </CardContent>
        <div className="absolute -top-2 -left-2 h-6 w-6 rounded-full border-4 border-black bg-white"></div>
        <div className="absolute -right-2 -bottom-2 h-6 w-6 rounded-full border-4 border-black bg-white"></div>
      </Card>
    </div>
  );
}
