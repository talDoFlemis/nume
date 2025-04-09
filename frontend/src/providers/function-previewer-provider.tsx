import { createContext, useContext, useState } from "react";

type FunctionPreviewerState = {
  initialValue: string;
  setInitialValue: (value: string) => void;
  rawText: string;
  setRawText: (value: string) => void;
  rawLatex: string;
  setRawLatex: (value: string) => void;
};

const initialState: FunctionPreviewerState = {
  initialValue: "",
  setInitialValue: () => null,
  rawText: "",
  setRawText: () => null,
  rawLatex: "",
  setRawLatex: () => null,
};

const FunctionPreviewerContext = createContext<FunctionPreviewerState | null>(
  initialState,
);

export function FunctionPreviewerProvider({
  children,
}: {
  children: React.ReactNode;
}) {
  const [initialValue, setInitialValue] = useState<string>("");
  const [rawText, setRawText] = useState<string>("");
  const [rawLatex, setRawLatex] = useState<string>("");

  const value = {
    initialValue,
    setInitialValue,
    rawText,
    setRawText,
    rawLatex,
    setRawLatex,
  };

  return (
    <FunctionPreviewerContext.Provider value={value}>
      {children}
    </FunctionPreviewerContext.Provider>
  );
}

export function useFunctionPreviewer() {
  const context = useContext(FunctionPreviewerContext);
  if (context === null) {
    throw new Error(
      "useFunctionPreviewer must be used within a FunctionPreviewerProvider",
    );
  }
  return context;
}
