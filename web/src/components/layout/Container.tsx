import { cn } from "../../lib/utils";

interface ContainerProps {
  children: React.ReactNode;
  className?: string;
}

export function Container({ children, className }: ContainerProps) {
  return (
    <div className={cn("max-w-4xl mx-auto p-6", className)}>
      {children}
    </div>
  );
}