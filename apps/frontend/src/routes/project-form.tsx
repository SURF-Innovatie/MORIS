import { useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { useNavigate, useParams } from 'react-router-dom';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { CalendarIcon, Loader2 } from 'lucide-react';
import { format } from 'date-fns';

import { Button } from '../components/ui/button';
import {
    Form,
    FormControl,
    FormField,
    FormItem,
    FormLabel,
    FormMessage,
} from '../components/ui/form';
import { Input } from '../components/ui/input';
import { Textarea } from '../components/ui/textarea';
import {
    Popover,
    PopoverContent,
    PopoverTrigger,
} from '../components/ui/popover';
import { Calendar } from '../components/ui/calendar';
import { cn } from '../lib/utils';
import { useToast } from '../hooks/use-toast';
import {
    useGetProjectsId,
    usePostProjects,
    usePutProjectsId,
} from '../api/generated-orval/moris';

const formSchema = z.object({
    title: z.string().min(1, 'Title is required'),
    description: z.string().min(1, 'Description is required'),
    startDate: z.date({
        required_error: 'Start date is required',
    }),
    endDate: z.date({
        required_error: 'End date is required',
    }),
    organisationID: z.string().uuid('Invalid organisation ID'),
});

export default function ProjectFormRoute() {
    const { id } = useParams();
    const navigate = useNavigate();
    const { toast } = useToast();
    const isEditing = !!id;

    const { data: project, isLoading: isLoadingProject } = useGetProjectsId(id!, {
        query: {
            enabled: isEditing,
        },
    });

    const { mutateAsync: createProject, isPending: isCreating } = usePostProjects();
    const { mutateAsync: updateProject, isPending: isUpdating } = usePutProjectsId();

    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            title: '',
            description: '',
            // TODO: This should be dynamic or selected from a list
            organisationID: '00000000-0000-0000-0000-000000000000',
        },
    });

    useEffect(() => {
        if (project) {
            form.reset({
                title: project.title || '',
                description: project.description || '',
                startDate: project.startDate ? new Date(project.startDate) : undefined,
                endDate: project.endDate ? new Date(project.endDate) : undefined,
                organisationID: project.organization?.id || '00000000-0000-0000-0000-000000000000',
            });
        }
    }, [project, form]);

    async function onSubmit(values: z.infer<typeof formSchema>) {
        try {
            if (isEditing) {
                await updateProject({
                    id: id!,
                    data: {
                        title: values.title,
                        description: values.description,
                        startDate: values.startDate.toISOString(),
                        endDate: values.endDate.toISOString(),
                        organisationID: values.organisationID,
                    },
                });
                toast({
                    title: 'Project updated',
                    description: 'The project has been successfully updated.',
                });
            } else {
                await createProject({
                    data: {
                        title: values.title,
                        description: values.description,
                        startDate: values.startDate.toISOString(),
                        endDate: values.endDate.toISOString(),
                        organisationID: values.organisationID,
                    },
                });
                toast({
                    title: 'Project created',
                    description: 'The new project has been successfully created.',
                });
            }
            navigate('/dashboard');
        } catch (error) {
            toast({
                variant: 'destructive',
                title: 'Error',
                description: 'Something went wrong. Please try again.',
            });
        }
    }

    if (isEditing && isLoadingProject) {
        return (
            <div className="flex h-full items-center justify-center">
                <Loader2 className="h-8 w-8 animate-spin text-muted-foreground" />
            </div>
        );
    }

    return (
        <div className="mx-auto max-w-2xl py-8">
            <div className="mb-8">
                <h1 className="text-3xl font-bold tracking-tight">
                    {isEditing ? 'Edit Project' : 'Create Project'}
                </h1>
                <p className="text-muted-foreground">
                    {isEditing
                        ? 'Update the project details below.'
                        : 'Fill in the details to start a new project.'}
                </p>
            </div>

            <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)} className="space-y-8">
                    <FormField
                        control={form.control}
                        name="title"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Title</FormLabel>
                                <FormControl>
                                    <Input placeholder="Project title" {...field} />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />

                    <FormField
                        control={form.control}
                        name="description"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Description</FormLabel>
                                <FormControl>
                                    <Textarea
                                        placeholder="Describe the project..."
                                        className="resize-none"
                                        {...field}
                                    />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />

                    <div className="grid gap-4 sm:grid-cols-2">
                        <FormField
                            control={form.control}
                            name="startDate"
                            render={({ field }) => (
                                <FormItem className="flex flex-col">
                                    <FormLabel>Start Date</FormLabel>
                                    <Popover>
                                        <PopoverTrigger asChild>
                                            <FormControl>
                                                <Button
                                                    variant={'outline'}
                                                    className={cn(
                                                        'w-full pl-3 text-left font-normal',
                                                        !field.value && 'text-muted-foreground'
                                                    )}
                                                >
                                                    {field.value ? (
                                                        format(field.value, 'PPP')
                                                    ) : (
                                                        <span>Pick a date</span>
                                                    )}
                                                    <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                                                </Button>
                                            </FormControl>
                                        </PopoverTrigger>
                                        <PopoverContent className="w-auto p-0" align="start">
                                            <Calendar
                                                mode="single"
                                                selected={field.value}
                                                onSelect={field.onChange}
                                                disabled={(date) =>
                                                    date < new Date('1900-01-01')
                                                }
                                                initialFocus
                                            />
                                        </PopoverContent>
                                    </Popover>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />

                        <FormField
                            control={form.control}
                            name="endDate"
                            render={({ field }) => (
                                <FormItem className="flex flex-col">
                                    <FormLabel>End Date</FormLabel>
                                    <Popover>
                                        <PopoverTrigger asChild>
                                            <FormControl>
                                                <Button
                                                    variant={'outline'}
                                                    className={cn(
                                                        'w-full pl-3 text-left font-normal',
                                                        !field.value && 'text-muted-foreground'
                                                    )}
                                                >
                                                    {field.value ? (
                                                        format(field.value, 'PPP')
                                                    ) : (
                                                        <span>Pick a date</span>
                                                    )}
                                                    <CalendarIcon className="ml-auto h-4 w-4 opacity-50" />
                                                </Button>
                                            </FormControl>
                                        </PopoverTrigger>
                                        <PopoverContent className="w-auto p-0" align="start">
                                            <Calendar
                                                mode="single"
                                                selected={field.value}
                                                onSelect={field.onChange}
                                                disabled={(date) =>
                                                    date < new Date('1900-01-01')
                                                }
                                                initialFocus
                                            />
                                        </PopoverContent>
                                    </Popover>
                                    <FormMessage />
                                </FormItem>
                            )}
                        />
                    </div>

                    <FormField
                        control={form.control}
                        name="organisationID"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>Organisation ID</FormLabel>
                                <FormControl>
                                    <Input placeholder="UUID" {...field} />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />

                    <div className="flex justify-end gap-4">
                        <Button
                            type="button"
                            variant="ghost"
                            onClick={() => navigate('/dashboard')}
                        >
                            Cancel
                        </Button>
                        <Button type="submit" disabled={isCreating || isUpdating}>
                            {isCreating || isUpdating ? (
                                <>
                                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                                    Saving...
                                </>
                            ) : (
                                'Save Project'
                            )}
                        </Button>
                    </div>
                </form>
            </Form>
        </div>
    );
}
