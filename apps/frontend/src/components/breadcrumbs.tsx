import { useLocation, Link } from "react-router-dom";
import { ChevronRight, Home } from "lucide-react";

export const Breadcrumbs = () => {
    const location = useLocation();
    const pathnames = location.pathname.split("/").filter((x) => x);

    // Don't show breadcrumbs on the dashboard root if you don't want to,
    // or always show "Home > Dashboard".
    // For now, let's map "dashboard" to "Home" visually or just keep it as is.

    return (
        <nav aria-label="Breadcrumb" className="mb-4 flex items-center text-sm text-muted-foreground">
            <Link
                to="/dashboard"
                className="flex items-center hover:text-foreground transition-colors"
            >
                <Home className="h-4 w-4" />
            </Link>
            {pathnames.length > 0 && (
                <ChevronRight className="mx-2 h-4 w-4 text-muted-foreground/50" />
            )}
            {pathnames.map((value, index) => {
                const to = `/${pathnames.slice(0, index + 1).join("/")}`;
                const isLast = index === pathnames.length - 1;
                const label = value.charAt(0).toUpperCase() + value.slice(1);

                return (
                    <div key={to} className="flex items-center">
                        {isLast ? (
                            <span className="font-medium text-foreground">{label}</span>
                        ) : (
                            <Link to={to} className="hover:text-foreground transition-colors">
                                {label}
                            </Link>
                        )}
                        {!isLast && (
                            <ChevronRight className="mx-2 h-4 w-4 text-muted-foreground/50" />
                        )}
                    </div>
                );
            })}
        </nav>
    );
};
