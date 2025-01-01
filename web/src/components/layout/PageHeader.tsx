interface PageHeaderProps {
    title: string;
    description?: string;
  }
  
  export function PageHeader({ title, description }: PageHeaderProps) {
    return (
      <div className="mb-6">
        <h1 className="text-3xl font-bold text-foreground">{title}</h1>
        {description && (
          <p className="mt-2 text-muted-foreground">{description}</p>
        )}
      </div>
    );
  }