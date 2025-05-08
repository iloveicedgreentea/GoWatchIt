import { cn } from "../../lib/utils";
import { HTMLAttributes } from "react";

interface FormProps extends HTMLAttributes<HTMLFormElement> {
  children: React.ReactNode;
}

export function Form({ children, className, ...props }: FormProps) {
  return (
    <form {...props} className={cn("space-y-6", className)}>
      {children}
    </form>
  );
}

interface FloatingButtonProps extends HTMLAttributes<HTMLDivElement> {
  children: React.ReactNode;
}

export function FloatingButton({ children, className, ...props }: FloatingButtonProps) {
  return (
    <div {...props} className={cn("fixed bottom-6 right-6", className)}>
      {children}
    </div>
  );
}

// Extend the existing button component with our common styling
interface SaveButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  children: React.ReactNode;
}

export function SaveButton({ children, className, ...props }: SaveButtonProps) {
  return (
    <button
      type="submit"
      className={cn(
        "px-6 py-3 bg-primary text-primary-foreground rounded-md hover:bg-primary/90 transition-colors",
        className
      )}
      {...props}
    >
      {children}
    </button>
  );
}