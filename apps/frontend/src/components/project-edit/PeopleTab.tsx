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

import { PersonResponse } from "@api/model";

interface PeopleTabProps {
  projectId: string;
  people: PersonResponse[];
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
                <Avatar className="h-10 w-10 border">
                  <AvatarImage src={undefined} />
                  <AvatarFallback className="font-semibold text-primary">
                    {(person.name || "Unknown")
                      .split(" ")
                      .map((n) => n[0])
                      .join("")
                      .toUpperCase()
                      .slice(0, 2)}
                  </AvatarFallback>
                </Avatar>
                <div>
                  <div className="flex items-center gap-2">
                    <p className="font-semibold leading-none">{person.name || "Unknown"}</p>
                    {adminId === person.id && (
                      <Crown className="h-3.5 w-3.5 text-yellow-500 fill-yellow-500" />
                    )}
                    <Badge variant="secondary" className="text-[10px] h-5 px-1.5 font-normal">
                      Collaborator
                    </Badge>
                  </div>
                  <p className="text-sm text-muted-foreground mt-1">
                    {person.email}
                  </p>
                </div>
              </div>
              <div className="flex items-center gap-2">
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="ghost" size="icon" className="h-8 w-8">
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
