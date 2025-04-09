import { addStyles, EditableMathField } from "react-mathquill";
import { CardContent, CardTitle } from "../ui/card";
import { Button } from "../ui/button";
import { cn } from "@/lib/utils";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "../ui/tabs";
import { Divide, Minus, Plus, X } from "lucide-react";
import { ReactNode } from "react";
addStyles();

type FunctionPreviewerButtonProps = {
  children: ReactNode;
  // updater: (selectedText: string) => void;
} & React.ComponentProps<"button">;

const FunctionPreviewerButton = ({
  children,
  // updater,
  className,
  ...props
}: FunctionPreviewerButtonProps) => {
  return (
    <Button
      variant="outline"
      className={cn(
        "border-foreground h-10 cursor-pointer border-2 font-bold",
        className,
      )}
      onClick={(e) => {
        e.preventDefault();
        // updater(children);
      }}
      {...props}
    >
      {children}
    </Button>
  );
};

const FunctionPreviewer = () => {
  return (
    <div className="text-card-foreground group border-card-foreground bg-card relative flex flex-col gap-6 overflow-hidden rounded-lg border-4 py-4 shadow-[8px_8px_0px_0px] transition-all duration-200">
      <CardTitle className="bg-amber-300 px-4">
        <EditableMathField
          className="border-none"
          latex={"\\frac{1}{\\sqrt{2}}\\cdot 2"}
          onChange={(mathField) => {
            console.log(mathField.latex());
          }}
        />
      </CardTitle>
      <div className="border-b-4 border-black bg-gradient-to-r from-blue-50 to-purple-50 p-2"></div>
      <CardContent>
        <Tabs defaultValue="basic">
          <TabsList className="border-card-foreground bg-card mb-2 grid w-full grid-cols-3 border-2">
            <TabsTrigger
              value="basic"
              className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground text-xs"
            >
              Basic
            </TabsTrigger>
            <TabsTrigger
              value="functions"
              className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground text-xs"
            >
              Functions
            </TabsTrigger>
            <TabsTrigger
              value="advanced"
              className="data-[state=active]:bg-primary data-[state=active]:text-primary-foreground text-xs"
            >
              Advanced
            </TabsTrigger>
          </TabsList>
          <TabsContent value="basic">
            <div className="grid grid-cols-4 gap-1">
              <FunctionPreviewerButton>
                <Plus className="h-4 w-4" />
              </FunctionPreviewerButton>
              <FunctionPreviewerButton>
                <Minus className="h-4 w-4" />
              </FunctionPreviewerButton>
              <FunctionPreviewerButton>
                <X className="h-4 w-4" />
              </FunctionPreviewerButton>
              <FunctionPreviewerButton>
                <Divide className="h-4 w-4" />
              </FunctionPreviewerButton>
              <FunctionPreviewerButton>
                x<sup>y</sup>
              </FunctionPreviewerButton>
              <FunctionPreviewerButton>√</FunctionPreviewerButton>
              <FunctionPreviewerButton>(</FunctionPreviewerButton>
              <FunctionPreviewerButton>)</FunctionPreviewerButton>
              <FunctionPreviewerButton>x</FunctionPreviewerButton>
              <FunctionPreviewerButton>π</FunctionPreviewerButton>
              <FunctionPreviewerButton>e</FunctionPreviewerButton>
              <FunctionPreviewerButton>=</FunctionPreviewerButton>
            </div>
          </TabsContent>
        </Tabs>
      </CardContent>
    </div>
  );
};

export default FunctionPreviewer;
