"use client";

import { useState } from "react";
import { LineChart, Calculator } from "lucide-react";

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

interface MethodPlaygroundProps {}

const methods = ["Bisection Method", "Newton-Raphson Method", "Secant Method"];

export function RootFindingTab({}: MethodPlaygroundProps) {
  const [selectedMethod, setSelectedMethod] = useState(methods[0]);
  const [functionInput, setFunctionInput] = useState("x^2 - 4");
  const [paramA, setParamA] = useState("0");
  const [paramB, setParamB] = useState("3");
  const [iterations, setIterations] = useState("10");
  const [result, setResult] = useState<string | null>(null);
  const [isCalculating, setIsCalculating] = useState(false);

  const handleCalculate = () => {
    setIsCalculating(true);
    // Simulate calculation delay
    setTimeout(() => {
      setResult(
        `Method: ${selectedMethod}\nFunction: ${functionInput}\nParameters: [${paramA}, ${paramB}]\nIterations: ${iterations}\n\nResult would be displayed here with step-by-step calculations and visualization.`,
      );
      setIsCalculating(false);
    }, 1000);
  };

  return (
    <div className="mt-6 grid gap-12 py-6 lg:grid-cols-2">
      <Card className="py-6 transition-none">
        <div className="absolute -top-20 -left-20 h-40 w-40 rounded-full bg-blue-300 opacity-10"></div>
        <CardContent>
          <div className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="method" className="font-comic">
                Method
              </Label>
              <Select value={selectedMethod} onValueChange={setSelectedMethod}>
                <SelectTrigger className="border-2 border-black">
                  <SelectValue placeholder="Select method" />
                </SelectTrigger>
                <SelectContent className="group border-card-foreground bg-card relative overflow-hidden rounded-lg border-4 p-2 shadow-[8px_8px_0px_0px]">
                  {methods.map((method) => (
                    <SelectItem key={method} value={method}>
                      {method}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="function" className="font-comic">
                Function f(x)
              </Label>
              <Input
                id="function"
                value={functionInput}
                onChange={(e) => setFunctionInput(e.target.value)}
                placeholder="e.g. x^2 - 4"
                className="border-2 border-black"
              />
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div className="space-y-2">
                <Label htmlFor="paramA" className="font-comic">
                  Parameter A
                </Label>
                <Input
                  id="paramA"
                  value={paramA}
                  onChange={(e) => setParamA(e.target.value)}
                  placeholder="e.g. 0"
                  className="border-2 border-black"
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="paramB" className="font-comic">
                  Parameter B
                </Label>
                <Input
                  id="paramB"
                  value={paramB}
                  onChange={(e) => setParamB(e.target.value)}
                  placeholder="e.g. 1"
                  className="border-2 border-black"
                />
              </div>
            </div>

            <div className="space-y-2">
              <Label htmlFor="iterations" className="font-comic">
                Iterations
              </Label>
              <Input
                id="iterations"
                value={iterations}
                onChange={(e) => setIterations(e.target.value)}
                placeholder="e.g. 10"
                className="border-2 border-black"
              />
            </div>

            <Button
              onClick={handleCalculate}
              variant={"cartoon"}
              className="w-full cursor-pointer"
              disabled={isCalculating}
            >
              {isCalculating ? "Calculating..." : "Calculate"}
              <Calculator className="ml-2 h-4 w-4" />
            </Button>
          </div>
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
                value={result || "Results will appear here after calculation"}
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
