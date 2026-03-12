'use client';

import { useEffect, useState } from 'react';
import { MainLayout } from '@/src/presentation/components/layout/main-layout';
import { ProtectedRoute } from '@/src/presentation/components/layout/protected-route';
import { ProtectedModule } from '@/src/presentation/components/layout/protected-feature';
import { Button } from '@/src/presentation/components/ui/button';
import { Input } from '@/src/presentation/components/ui/input';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/src/presentation/components/ui/card';
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/src/presentation/components/ui/dialog';
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from '@/src/presentation/components/ui/alert-dialog';
import { Label } from '@/src/presentation/components/ui/label';
import { Textarea } from '@/src/presentation/components/ui/textarea';
import { Alert, AlertDescription } from '@/src/presentation/components/ui/alert';
import { Badge } from '@/src/presentation/components/ui/badge';
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from '@/src/presentation/components/ui/table';
import { Plus, Edit, Trash2, Shield, Loader2, CheckCircle, Search, Users, Key, Settings } from 'lucide-react';
import { rbacService, type Role, type RoleWithPermissions, type Permission } from '@/src/domain/services/rbac.service';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/src/presentation/components/ui/tabs';
import { ScrollArea } from '@/src/presentation/components/ui/scroll-area';

export default function RolesPage() {
    const [roles, setRoles] = useState<Role[]>([]);
    const [permissions, setPermissions] = useState<Permission[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState('');
    const [searchQuery, setSearchQuery] = useState('');

    // Dialog states
    const [isCreateDialogOpen, setIsCreateDialogOpen] = useState(false);
    const [isEditDialogOpen, setIsEditDialogOpen] = useState(false);
    const [isDeleteDialogOpen, setIsDeleteDialogOpen] = useState(false);
    const [isPermissionsDialogOpen, setIsPermissionsDialogOpen] = useState(false);
    const [selectedRole, setSelectedRole] = useState<RoleWithPermissions | null>(null);
    const [selectedPermissions, setSelectedPermissions] = useState<number[]>([]);
    const [isSubmitting, setIsSubmitting] = useState(false);

    // Form states
    const [formData, setFormData] = useState({
        name: '',
        description: '',
    });

    useEffect(() => {
        loadData();
    }, []);

    const loadData = async () => {
        try {
            setLoading(true);
            const [rolesData, permissionsData] = await Promise.all([
                rbacService.getAllRoles(),
                rbacService.getAllPermissions(),
            ]);

            setRoles(rolesData);
            setPermissions(permissionsData);
            setLoading(false);
        } catch (err: any) {
            setError(err.response?.data?.message || 'Failed to load data');
            setLoading(false);
        }
    };

    const resetForm = () => {
        setFormData({ name: '', description: '' });
        setSelectedRole(null);
        setSelectedPermissions([]);
        setError('');
    };

    const handleCreateRole = async () => {
        try {
            setIsSubmitting(true);
            await rbacService.createRole({
                name: formData.name,
                description: formData.description || undefined,
            });

            setIsCreateDialogOpen(false);
            resetForm();
            await loadData();
        } catch (err: any) {
            setError(err.response?.data?.message || 'Failed to create role');
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleUpdateRole = async () => {
        if (!selectedRole) return;

        try {
            setIsSubmitting(true);
            await rbacService.updateRole(selectedRole.id, {
                name: formData.name,
                description: formData.description || undefined,
            });

            setIsEditDialogOpen(false);
            resetForm();
            await loadData();
        } catch (err: any) {
            setError(err.response?.data?.message || 'Failed to update role');
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleDeleteRole = async () => {
        if (!selectedRole) return;

        try {
            setIsSubmitting(true);
            await rbacService.deleteRole(selectedRole.id);

            setIsDeleteDialogOpen(false);
            setSelectedRole(null);
            await loadData();
        } catch (err: any) {
            setError(err.response?.data?.message || 'Failed to delete role');
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleAssignPermissions = async () => {
        if (!selectedRole) return;

        try {
            setIsSubmitting(true);
            await rbacService.assignPermissionsToRole(selectedRole.id, {
                permission_ids: selectedPermissions,
            });

            setIsPermissionsDialogOpen(false);
            setSelectedRole(null);
            await loadData();
        } catch (err: any) {
            setError(err.response?.data?.message || 'Failed to assign permissions');
        } finally {
            setIsSubmitting(false);
        }
    };

    const openEditDialog = (role: Role) => {
        setFormData({ name: role.name, description: role.description || '' });
        setSelectedRole(role as RoleWithPermissions);
        setIsEditDialogOpen(true);
    };

    const openDeleteDialog = (role: Role) => {
        setSelectedRole(role as RoleWithPermissions);
        setIsDeleteDialogOpen(true);
    };

    const openPermissionsDialog = async (role: Role) => {
        try {
            const roleWithPerms = await rbacService.getRoleById(role.id);
            setSelectedRole(roleWithPerms);
            setSelectedPermissions(roleWithPerms.permissions.map(p => p.id));
            setIsPermissionsDialogOpen(true);
        } catch (err: any) {
            setError(err.response?.data?.message || 'Failed to load role permissions');
        }
    };

    // Filter roles by search
    const filteredRoles = roles.filter(role =>
        role.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
        (role.description && role.description.toLowerCase().includes(searchQuery.toLowerCase()))
    );

    // Group permissions by module
    const permissionsByModule = permissions.reduce((acc, perm) => {
        if (!acc[perm.module]) {
            acc[perm.module] = [];
        }
        acc[perm.module].push(perm);
        return acc;
    }, {} as Record<string, Permission[]>);

    // Count permissions for a role (for display in table)
    const getPermissionCount = async (roleId: number): Promise<number> => {
        try {
            const roleWithPerms = await rbacService.getRoleById(roleId);
            return roleWithPerms.permissions.length;
        } catch {
            return 0;
        }
    };

    // Toggle all permissions in a module
    const toggleModulePermissions = (modulePermissions: Permission[], isSelected: boolean) => {
        const modulePermIds = modulePermissions.map(p => p.id);
        if (isSelected) {
            // Remove all permissions from this module
            setSelectedPermissions(prev => prev.filter(id => !modulePermIds.includes(id)));
        } else {
            // Add all permissions from this module
            setSelectedPermissions(prev => [...new Set([...prev, ...modulePermIds])]);
        }
    };

    // Check if all permissions in a module are selected
    const isModuleFullySelected = (modulePermissions: Permission[]) => {
        return modulePermissions.every(perm => selectedPermissions.includes(perm.id));
    };

    return (
        <ProtectedRoute>
            <ProtectedModule requiredRole={['Super Admin', 'Admin']}>
                <MainLayout>
                    <Card>
                        <CardHeader>
                            <div className="flex flex-col md:flex-row justify-between items-start gap-4">
                                <div>
                                    <CardTitle className="text-3xl font-bold flex items-center gap-2">
                                        Role Management
                                    </CardTitle>
                                    <CardDescription className="mt-2">
                                        Create and manage roles with specific permissions for access control
                                    </CardDescription>
                                </div>
                                <Button onClick={() => { resetForm(); setIsCreateDialogOpen(true); }}>
                                    <Plus className="mr-2 h-4 w-4" />
                                    Create Role
                                </Button>
                            </div>
                        </CardHeader>
                        <CardContent>
                            {error && (
                                <Alert variant="destructive" className="mb-4">
                                    <AlertDescription>{error}</AlertDescription>
                                </Alert>
                            )}

                            {/* Search */}
                            <div className="flex gap-4 mb-6">
                                <div className="relative flex-1 max-w-md">
                                    <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground h-4 w-4" />
                                    <Input
                                        placeholder="Search roles..."
                                        value={searchQuery}
                                        onChange={(e) => setSearchQuery(e.target.value)}
                                        className="pl-10"
                                    />
                                </div>
                            </div>

                            {/* Stats Cards */}
                            <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
                                <Card>
                                    <CardContent className="flex items-center gap-4 pt-6">
                                        <div className="p-3 rounded-full bg-blue-500/20">
                                            <Shield className="h-6 w-6 text-blue-500" />
                                        </div>
                                        <div>
                                            <p className="text-2xl font-bold">{roles.length}</p>
                                            <p className="text-sm text-muted-foreground">Total Roles</p>
                                        </div>
                                    </CardContent>
                                </Card>
                                <Card>
                                    <CardContent className="flex items-center gap-4 pt-6">
                                        <div className="p-3 rounded-full bg-green-500/20">
                                            <Key className="h-6 w-6 text-green-500" />
                                        </div>
                                        <div>
                                            <p className="text-2xl font-bold">{permissions.length}</p>
                                            <p className="text-sm text-muted-foreground">Total Permissions</p>
                                        </div>
                                    </CardContent>
                                </Card>
                                <Card>
                                    <CardContent className="flex items-center gap-4 pt-6">
                                        <div className="p-3 rounded-full bg-purple-500/20">
                                            <Settings className="h-6 w-6 text-purple-500" />
                                        </div>
                                        <div>
                                            <p className="text-2xl font-bold">{Object.keys(permissionsByModule).length}</p>
                                            <p className="text-sm text-muted-foreground">Modules</p>
                                        </div>
                                    </CardContent>
                                </Card>
                            </div>

                            {loading ? (
                                <div className="flex justify-center items-center h-64">
                                    <Loader2 className="h-8 w-8 animate-spin text-primary" />
                                </div>
                            ) : filteredRoles.length === 0 ? (
                                <div className="text-center py-12">
                                    <Shield className="h-16 w-16 mx-auto text-muted-foreground mb-4" />
                                    <h3 className="text-lg font-semibold mb-2">No Roles Found</h3>
                                    <p className="text-muted-foreground">
                                        {searchQuery ? 'Try a different search term' : 'Create your first role to get started'}
                                    </p>
                                </div>
                            ) : (
                                <div className="border rounded-lg overflow-hidden">
                                    <Table>
                                        <TableHeader>
                                            <TableRow>
                                                <TableHead>Role Name</TableHead>
                                                <TableHead>Description</TableHead>
                                                <TableHead>Type</TableHead>
                                                <TableHead>Created</TableHead>
                                                <TableHead className="text-right">Actions</TableHead>
                                            </TableRow>
                                        </TableHeader>
                                        <TableBody>
                                            {filteredRoles.map((role) => (
                                                <TableRow key={role.id}>
                                                    <TableCell>
                                                        <div className="flex items-center gap-2">
                                                            <div className="p-2 rounded-full bg-primary/10">
                                                                <Shield className="h-4 w-4 text-primary" />
                                                            </div>
                                                            <span className="font-medium">{role.name}</span>
                                                        </div>
                                                    </TableCell>
                                                    <TableCell className="text-muted-foreground">
                                                        {role.description || 'No description'}
                                                    </TableCell>
                                                    <TableCell>
                                                        {role.is_system ? (
                                                            <Badge variant="secondary">System</Badge>
                                                        ) : (
                                                            <Badge variant="outline">Custom</Badge>
                                                        )}
                                                    </TableCell>
                                                    <TableCell className="text-muted-foreground">
                                                        {new Date(role.created_at).toLocaleDateString()}
                                                    </TableCell>
                                                    <TableCell className="text-right">
                                                        <div className="flex justify-end gap-2">
                                                            <Button
                                                                size="sm"
                                                                variant="outline"
                                                                onClick={() => openPermissionsDialog(role)}
                                                                title="Manage Permissions"
                                                            >
                                                                <Key className="h-4 w-4" />
                                                            </Button>
                                                            <Button
                                                                size="sm"
                                                                variant="outline"
                                                                onClick={() => openEditDialog(role)}
                                                                disabled={role.is_system}
                                                                title={role.is_system ? 'System roles cannot be edited' : 'Edit Role'}
                                                            >
                                                                <Edit className="h-4 w-4" />
                                                            </Button>
                                                            <Button
                                                                size="sm"
                                                                variant="outline"
                                                                onClick={() => openDeleteDialog(role)}
                                                                disabled={role.is_system}
                                                                title={role.is_system ? 'System roles cannot be deleted' : 'Delete Role'}
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
                            )}
                        </CardContent>
                    </Card>

                    {/* Create Role Dialog */}
                    <Dialog open={isCreateDialogOpen} onOpenChange={setIsCreateDialogOpen}>
                        <DialogContent>
                            <DialogHeader>
                                <DialogTitle>Create New Role</DialogTitle>
                                <DialogDescription>
                                    Create a new role to define user access levels
                                </DialogDescription>
                            </DialogHeader>

                            <div className="space-y-4">
                                <div className="space-y-2">
                                    <Label htmlFor="create-name">Role Name</Label>
                                    <Input
                                        id="create-name"
                                        value={formData.name}
                                        onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
                                        placeholder="e.g., Editor, Viewer"
                                    />
                                </div>

                                <div className="space-y-2">
                                    <Label htmlFor="create-description">Description</Label>
                                    <Textarea
                                        id="create-description"
                                        value={formData.description}
                                        onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
                                        placeholder="Describe what this role can do..."
                                        rows={3}
                                    />
                                </div>
                            </div>

                            <DialogFooter>
                                <Button variant="outline" onClick={() => setIsCreateDialogOpen(false)}>
                                    Cancel
                                </Button>
                                <Button onClick={handleCreateRole} disabled={!formData.name.trim() || isSubmitting}>
                                    {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                                    Create
                                </Button>
                            </DialogFooter>
                        </DialogContent>
                    </Dialog>

                    {/* Edit Role Dialog */}
                    <Dialog open={isEditDialogOpen} onOpenChange={setIsEditDialogOpen}>
                        <DialogContent>
                            <DialogHeader>
                                <DialogTitle>Edit Role</DialogTitle>
                                <DialogDescription>
                                    Update role information
                                </DialogDescription>
                            </DialogHeader>

                            <div className="space-y-4">
                                <div className="space-y-2">
                                    <Label htmlFor="edit-name">Role Name</Label>
                                    <Input
                                        id="edit-name"
                                        value={formData.name}
                                        onChange={(e) => setFormData(prev => ({ ...prev, name: e.target.value }))}
                                    />
                                </div>

                                <div className="space-y-2">
                                    <Label htmlFor="edit-description">Description</Label>
                                    <Textarea
                                        id="edit-description"
                                        value={formData.description}
                                        onChange={(e) => setFormData(prev => ({ ...prev, description: e.target.value }))}
                                        rows={3}
                                    />
                                </div>
                            </div>

                            <DialogFooter>
                                <Button variant="outline" onClick={() => setIsEditDialogOpen(false)}>
                                    Cancel
                                </Button>
                                <Button onClick={handleUpdateRole} disabled={!formData.name.trim() || isSubmitting}>
                                    {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                                    Save Changes
                                </Button>
                            </DialogFooter>
                        </DialogContent>
                    </Dialog>

                    {/* Delete Role Dialog */}
                    <AlertDialog open={isDeleteDialogOpen} onOpenChange={setIsDeleteDialogOpen}>
                        <AlertDialogContent>
                            <AlertDialogHeader>
                                <AlertDialogTitle>Delete Role</AlertDialogTitle>
                                <AlertDialogDescription>
                                    Are you sure you want to delete the role &quot;{selectedRole?.name}&quot;?
                                    This action cannot be undone. Users with this role will lose their assigned permissions.
                                </AlertDialogDescription>
                            </AlertDialogHeader>
                            <AlertDialogFooter>
                                <AlertDialogCancel>Cancel</AlertDialogCancel>
                                <AlertDialogAction
                                    onClick={handleDeleteRole}
                                    className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                                >
                                    {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                                    Delete
                                </AlertDialogAction>
                            </AlertDialogFooter>
                        </AlertDialogContent>
                    </AlertDialog>

                    {/* Manage Permissions Dialog */}
                    <Dialog open={isPermissionsDialogOpen} onOpenChange={setIsPermissionsDialogOpen}>
                        <DialogContent className="max-w-4xl max-h-[85vh]">
                            <DialogHeader>
                                <DialogTitle className="flex items-center gap-2">
                                    <Key className="h-5 w-5 text-primary" />
                                    Manage Permissions
                                </DialogTitle>
                                <DialogDescription>
                                    Assign permissions to {selectedRole?.name}.
                                    Selected: {selectedPermissions.length} of {permissions.length} permissions
                                </DialogDescription>
                            </DialogHeader>

                            <div className="flex-1 overflow-hidden">
                                {Object.keys(permissionsByModule).length > 0 ? (
                                    <Tabs defaultValue={Object.keys(permissionsByModule)[0]} className="h-full">
                                        <TabsList className="w-full justify-start flex-wrap h-auto gap-1 p-1 bg-muted/50">
                                            {Object.keys(permissionsByModule).map((module) => (
                                                <TabsTrigger
                                                    key={module}
                                                    value={module}
                                                    className="capitalize data-[state=active]:bg-background"
                                                >
                                                    {module}
                                                    <Badge variant="secondary" className="ml-2 text-xs">
                                                        {permissionsByModule[module].filter(p => selectedPermissions.includes(p.id)).length}/
                                                        {permissionsByModule[module].length}
                                                    </Badge>
                                                </TabsTrigger>
                                            ))}
                                        </TabsList>

                                        {Object.entries(permissionsByModule).map(([module, perms]) => {
                                            const allSelected = isModuleFullySelected(perms);
                                            return (
                                                <TabsContent key={module} value={module} className="mt-4">
                                                    <div className="border rounded-lg">
                                                        {/* Module header with select all */}
                                                        <div className="flex items-center justify-between p-4 border-b bg-muted/30">
                                                            <div className="flex items-center gap-2">
                                                                <Settings className="h-4 w-4 text-muted-foreground" />
                                                                <span className="font-medium capitalize">{module} Module</span>
                                                            </div>
                                                            <Button
                                                                variant="outline"
                                                                size="sm"
                                                                onClick={() => toggleModulePermissions(perms, allSelected)}
                                                            >
                                                                {allSelected ? 'Deselect All' : 'Select All'}
                                                            </Button>
                                                        </div>

                                                        <ScrollArea className="h-[300px]">
                                                            <div className="p-2">
                                                                {perms.map((perm) => (
                                                                    <div
                                                                        key={perm.id}
                                                                        className={`flex items-start space-x-3 p-3 rounded-lg cursor-pointer transition-colors ${selectedPermissions.includes(perm.id)
                                                                            ? 'bg-primary/10 border border-primary/20'
                                                                            : 'hover:bg-muted/50'
                                                                            }`}
                                                                        onClick={() => {
                                                                            if (selectedPermissions.includes(perm.id)) {
                                                                                setSelectedPermissions(prev => prev.filter(id => id !== perm.id));
                                                                            } else {
                                                                                setSelectedPermissions(prev => [...prev, perm.id]);
                                                                            }
                                                                        }}
                                                                    >
                                                                        <input
                                                                            type="checkbox"
                                                                            checked={selectedPermissions.includes(perm.id)}
                                                                            onChange={() => { }}
                                                                            className="h-5 w-5 rounded border-gray-300 mt-0.5"
                                                                        />
                                                                        <div className="flex-1">
                                                                            <div className="font-medium">{perm.name}</div>
                                                                            {perm.description && (
                                                                                <div className="text-sm text-muted-foreground mt-1">
                                                                                    {perm.description}
                                                                                </div>
                                                                            )}
                                                                        </div>
                                                                        {selectedPermissions.includes(perm.id) && (
                                                                            <CheckCircle className="h-5 w-5 text-primary" />
                                                                        )}
                                                                    </div>
                                                                ))}
                                                            </div>
                                                        </ScrollArea>
                                                    </div>
                                                </TabsContent>
                                            );
                                        })}
                                    </Tabs>
                                ) : (
                                    <div className="text-center py-8">
                                        <Key className="h-12 w-12 mx-auto text-muted-foreground mb-4" />
                                        <p className="text-muted-foreground">No permissions available</p>
                                    </div>
                                )}
                            </div>

                            <DialogFooter className="mt-4">
                                <Button variant="outline" onClick={() => setIsPermissionsDialogOpen(false)}>
                                    Cancel
                                </Button>
                                <Button onClick={handleAssignPermissions} disabled={isSubmitting}>
                                    {isSubmitting && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                                    Save Permissions ({selectedPermissions.length})
                                </Button>
                            </DialogFooter>
                        </DialogContent>
                    </Dialog>
                </MainLayout>
            </ProtectedModule>
        </ProtectedRoute>
    );
}
