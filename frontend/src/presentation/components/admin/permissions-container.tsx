'use client';

import { useEffect, useState } from 'react';
import { Button } from '@/src/presentation/components/ui/button';
import { Input } from '@/src/presentation/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/src/presentation/components/ui/card';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/src/presentation/components/ui/dialog';
import { Label } from '@/src/presentation/components/ui/label';
import { Textarea } from '@/src/presentation/components/ui/textarea';
import { Alert, AlertDescription } from '@/src/presentation/components/ui/alert';
import { Badge } from '@/src/presentation/components/ui/badge';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/src/presentation/components/ui/table';
import { Plus, Edit, Trash2, Key, Loader2, Search } from 'lucide-react';
import { rbacService, type Permission, type PermissionsByModule } from '@/src/domain/services/rbac.service';

export function PermissionsContainer() {
    const [permissionsByModule, setPermissionsByModule] = useState<PermissionsByModule[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [searchQuery, setSearchQuery] = useState('');

    // Dialog states
    const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
    const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
    const [selectedPermission, setSelectedPermission] = useState<Permission | null>(null);

    // Form states
    const [permissionName, setPermissionName] = useState('');
    const [permissionModule, setPermissionModule] = useState('');
    const [permissionDescription, setPermissionDescription] = useState('');

    useEffect(() => {
        loadData();
    }, []);

    const loadData = async () => {
        try {
            setLoading(true);
            const data = await rbacService.getPermissionsByModule();
            setPermissionsByModule(data);
            setLoading(false);
        } catch (err: any) {
            setError(err.response?.data?.message || 'Failed to load permissions');
            setLoading(false);
        }
    };

    const handleCreatePermission = async () => {
        try {
            await rbacService.createPermission({
                name: permissionName,
                module: permissionModule,
                description: permissionDescription || undefined,
            });

            setIsCreateDialogOpen(false);
            resetForm();
            await loadData();
        } catch (err: any) {
            setError(err.response?.data?.message || 'Failed to create permission');
        }
    };

    const handleEditPermission = async () => {
        if (!selectedPermission) return;

        try {
            await rbacService.updatePermission(selectedPermission.id, {
                name: permissionName,
                module: permissionModule,
                description: permissionDescription || undefined,
            });

            setIsEditDialogOpen(false);
            setSelectedPermission(null);
            resetForm();
            await loadData();
        } catch (err: any) {
            setError(err.response?.data?.message || 'Failed to update permission');
        }
    };

    const handleDeletePermission = async (id: number) => {
        if (!confirm('Are you sure you want to delete this permission? This action cannot be undone.')) {
            return;
        }

        try {
            await rbacService.deletePermission(id);
            await loadData();
        } catch (err: any) {
            setError(err.response?.data?.message || 'Failed to delete permission');
        }
    };

    const openEditDialog = (permission: Permission) => {
        setSelectedPermission(permission);
        setPermissionName(permission.name);
        setPermissionModule(permission.module);
        setPermissionDescription(permission.description || '');
        setIsEditDialogOpen(true);
    };

    const resetForm = () => {
        setPermissionName('');
        setPermissionModule('');
        setPermissionDescription('');
    };

    const filteredPermissions = permissionsByModule.map(module => ({
        ...module,
        permissions: module.permissions.filter(perm =>
            perm.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
            perm.description?.toLowerCase().includes(searchQuery.toLowerCase())
        )
    })).filter(module => module.permissions.length > 0);

    const totalPermissions = permissionsByModule.reduce((sum, module) => sum + module.permissions.length, 0);

    return (
        <>
            <Card>
                <CardHeader>
                    <div className="flex justify-between items-start">
                        <div>
                            <CardTitle className="text-3xl font-bold">Permission Management</CardTitle>
                            <CardDescription>
                                Manage system permissions and organize them by module
                            </CardDescription>
                        </div>
                        <Button onClick={() => setIsCreateDialogOpen(true)}>
                            <Plus className="mr-2 h-4 w-4" />
                            Create Permission
                        </Button>
                    </div>
                </CardHeader>
                <CardContent>
                    {error && (
                        <Alert variant="destructive" className="mb-4">
                            <AlertDescription>{error}</AlertDescription>
                        </Alert>
                    )}

                    {/* Stats */}
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
                        <Card>
                            <CardHeader className="pb-2">
                                <CardDescription>Total Permissions</CardDescription>
                                <CardTitle className="text-3xl">{totalPermissions}</CardTitle>
                            </CardHeader>
                        </Card>
                        <Card>
                            <CardHeader className="pb-2">
                                <CardDescription>Modules</CardDescription>
                                <CardTitle className="text-3xl">{permissionsByModule.length}</CardTitle>
                            </CardHeader>
                        </Card>
                        <Card>
                            <CardHeader className="pb-2">
                                <CardDescription>Search Results</CardDescription>
                                <CardTitle className="text-3xl">
                                    {filteredPermissions.reduce((sum, m) => sum + m.permissions.length, 0)}
                                </CardTitle>
                            </CardHeader>
                        </Card>
                    </div>

                    {/* Search */}
                    <div className="relative mb-6">
                        <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
                        <Input
                            placeholder="Search permissions..."
                            value={searchQuery}
                            onChange={(e) => setSearchQuery(e.target.value)}
                            className="pl-10"
                        />
                    </div>

                    {/* Permissions by Module */}
                    {loading ? (
                        <div className="flex justify-center items-center h-64">
                            <Loader2 className="h-8 w-8 animate-spin text-primary" />
                        </div>
                    ) : (
                        <div className="space-y-6">
                            {filteredPermissions.map((moduleGroup) => (
                                <Card key={moduleGroup.module}>
                                    <CardHeader>
                                        <div className="flex items-center justify-between">
                                            <div className="flex items-center gap-2">
                                                <Key className="h-5 w-5 text-primary" />
                                                <CardTitle className="text-xl capitalize">
                                                    {moduleGroup.module}
                                                </CardTitle>
                                                <Badge variant="secondary">
                                                    {moduleGroup.permissions.length} permissions
                                                </Badge>
                                            </div>
                                        </div>
                                    </CardHeader>
                                    <CardContent>
                                        <div className="border rounded-lg overflow-hidden">
                                            <Table>
                                                <TableHeader>
                                                    <TableRow>
                                                        <TableHead>Name</TableHead>
                                                        <TableHead>Description</TableHead>
                                                        <TableHead>Created</TableHead>
                                                        <TableHead className="text-right">Actions</TableHead>
                                                    </TableRow>
                                                </TableHeader>
                                                <TableBody>
                                                    {moduleGroup.permissions.map((permission) => (
                                                        <TableRow key={permission.id}>
                                                            <TableCell className="font-medium">
                                                                {permission.name}
                                                            </TableCell>
                                                            <TableCell className="text-muted-foreground">
                                                                {permission.description || '-'}
                                                            </TableCell>
                                                            <TableCell className="text-sm text-muted-foreground">
                                                                {new Date(permission.created_at).toLocaleDateString()}
                                                            </TableCell>
                                                            <TableCell className="text-right">
                                                                <div className="flex justify-end gap-2">
                                                                    <Button
                                                                        size="sm"
                                                                        variant="outline"
                                                                        onClick={() => openEditDialog(permission)}
                                                                    >
                                                                        <Edit className="h-4 w-4" />
                                                                    </Button>
                                                                    <Button
                                                                        size="sm"
                                                                        variant="outline"
                                                                        onClick={() => handleDeletePermission(permission.id)}
                                                                    >
                                                                        <Trash2 className="h-4 w-4" />
                                                                    </Button>
                                                                </div>
                                                            </TableCell>
                                                        </TableRow>
                                                    ))}
                                                </TableBody>
                                            </Table>
                                        </div>
                                    </CardContent>
                                </Card>
                            ))}

                            {filteredPermissions.length === 0 && !loading && (
                                <div className="text-center py-12">
                                    <Key className="h-16 w-16 mx-auto text-muted-foreground mb-4" />
                                    <h3 className="text-lg font-semibold mb-2">No Permissions Found</h3>
                                    <p className="text-muted-foreground">
                                        {searchQuery ? 'No permissions match your search.' : 'Create your first permission to get started.'}
                                    </p>
                                </div>
                            )}
                        </div>
                    )}
                </CardContent>
            </Card>

            {/* Create Permission Dialog */}
            <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Create New Permission</DialogTitle>
                        <DialogDescription>
                            Add a new permission to the system
                        </DialogDescription>
                    </DialogHeader>

                    <div className="space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="name">Permission Name</Label>
                            <Input
                                id="name"
                                value={permissionName}
                                onChange={(e) => setPermissionName(e.target.value)}
                                placeholder="e.g., users.create"
                            />
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="module">Module</Label>
                            <Input
                                id="module"
                                value={permissionModule}
                                onChange={(e) => setPermissionModule(e.target.value)}
                                placeholder="e.g., users"
                            />
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="description">Description</Label>
                            <Textarea
                                id="description"
                                value={permissionDescription}
                                onChange={(e) => setPermissionDescription(e.target.value)}
                                placeholder="Describe what this permission allows..."
                                rows={3}
                            />
                        </div>
                    </div>

                    <DialogFooter>
                        <Button variant="outline" onClick={() => setIsCreateDialogOpen(false)}>
                            Cancel
                        </Button>
                        <Button
                            onClick={handleCreatePermission}
                            disabled={!permissionName.trim() || !permissionModule.trim()}
                        >
                            Create Permission
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>

            {/* Edit Permission Dialog */}
            <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
                <DialogContent>
                    <DialogHeader>
                        <DialogTitle>Edit Permission</DialogTitle>
                        <DialogDescription>
                            Update permission details
                        </DialogDescription>
                    </DialogHeader>

                    <div className="space-y-4">
                        <div className="space-y-2">
                            <Label htmlFor="edit-name">Permission Name</Label>
                            <Input
                                id="edit-name"
                                value={permissionName}
                                onChange={(e) => setPermissionName(e.target.value)}
                                placeholder="e.g., users.create"
                            />
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="edit-module">Module</Label>
                            <Input
                                id="edit-module"
                                value={permissionModule}
                                onChange={(e) => setPermissionModule(e.target.value)}
                                placeholder="e.g., users"
                            />
                        </div>

                        <div className="space-y-2">
                            <Label htmlFor="edit-description">Description</Label>
                            <Textarea
                                id="edit-description"
                                value={permissionDescription}
                                onChange={(e) => setPermissionDescription(e.target.value)}
                                placeholder="Describe what this permission allows..."
                                rows={3}
                            />
                        </div>
                    </div>

                    <DialogFooter>
                        <Button variant="outline" onClick={() => setIsEditDialogOpen(false)}>
                            Cancel
                        </Button>
                        <Button
                            onClick={handleEditPermission}
                            disabled={!permissionName.trim() || !permissionModule.trim()}
                        >
                            Save Changes
                        </Button>
                    </DialogFooter>
                </DialogContent>
            </Dialog>
        </>
    );
}
