import { MoreHorizontal, Crown } from "lucide-react";

import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Badge } from "@/components/ui/badge";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { AddPersonDialog } from "./AddPersonDialog";

export interface Person {
  id: string;
  name: string;
  email: string;
  role?: string; // Role might be optional or coming from a different place
  avatar?: string | null;
}

interface PeopleTabProps {
  projectId: string;
  people: Person[];
  adminId?: string;
  onRefresh: () => void;
}

export function PeopleTab({
  projectId,
  people,
  adminId,
  onRefresh,
}: PeopleTabProps) {
  return (
    <Card>
      <CardHeader className="flex flex-row items-center justify-between">
        <div>
          <CardTitle>Team Members</CardTitle>
          <CardDescription>
            Manage who has access to this project.
          </CardDescription>
        </div>
        <AddPersonDialog projectId={projectId} onPersonAdded={onRefresh} />
      </CardHeader>
      <CardContent>
        <div className="space-y-4">
          {people.map((person) => (
            <div
              key={person.id}
              className="flex items-center justify-between rounded-lg border p-4 hover:bg-muted/50 transition-colors"
            >
              <div className="flex items-center gap-4">
                <Avatar>
                  <AvatarImage src={person.avatar || undefined} />
                  <AvatarFallback>
                    {person.name
                      .split(" ")
                      .map((n) => n[0])
                      .join("")}
                  </AvatarFallback>
                </Avatar>
                <div>
                  <div className="flex items-center gap-2">
                    <p className="font-medium leading-none">{person.name}</p>
                    {adminId === person.id && (
                      <Crown className="h-3 w-3 text-yellow-500" />
                    )}
                  </div>
                  <p className="text-sm text-muted-foreground">
                    {person.email}
                  </p>
                </div>
              </div>
              <div className="flex items-center gap-4">
                {person.role && <Badge variant="outline">{person.role}</Badge>}
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" size="icon">
                      <MoreHorizontal className="h-4 w-4" />
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent align="end">
                    <DropdownMenuLabel>Actions</DropdownMenuLabel>
                    <DropdownMenuItem>Edit Role</DropdownMenuItem>
                    <DropdownMenuSeparator />
                    <DropdownMenuItem className="text-destructive">
                      Remove
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              </div>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
