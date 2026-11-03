
import api from "./api";

export interface WorkOrderAssignment {
  id: number;
  work_order_id: number;
  technician_id: number;
  assigned_at: string;
  created_at: string;
  updated_at: string;
}

export interface WorkOrder {
  id: number;
  service_request_id: number;
  customer_id: number;
  title: string;
  description: string;
  service_address: string;
  status: string;
  scheduled_start: string | null;
  scheduled_end: string | null;
  vehicle_assignment?: WorkOrderVehicle | null;
  assignment?: WorkOrderAssignment | null;
  created_at: string;
  updated_at: string;
}

export interface WorkOrderVehicle {
  id: number;
  work_order_id: number;
  vehicle_id: number;
  assigned_at: string;
  created_at: string;
  updated_at: string;
}

interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

export const getWorkOrders = async (): Promise<WorkOrder[]> => {
  const response = await api.get<ApiResponse<WorkOrder[]>>("/work-orders");

  if (!response.data.success) {
    throw new Error(
      response.data.message || "Failed to retrieve work orders",
    );
  }

  return response.data.data;
};

export const getWorkOrderById = async (
  id: number,
): Promise<WorkOrder> => {
  const response = await api.get<ApiResponse<WorkOrder>>(
    `/work-orders/${id}`,
  );

  if (!response.data.success) {
    throw new Error(
      response.data.message || "Failed to retrieve work order",
    );
  }

  return response.data.data;
};

export const assignTechnician = async (
  workOrderId: number,
  technicianId: number,
): Promise<WorkOrderAssignment> => {
  const response = await api.post<ApiResponse<WorkOrderAssignment>>(
    `/work-orders/${workOrderId}/assign-technician`,
    { technician_id: technicianId },
  );

  if (!response.data.success) {
    throw new Error(
      response.data.message || "Failed to assign technician",
    );
  }

  return response.data.data;
};

export const reassignTechnician = async (
  workOrderId: number,
  technicianId: number,
): Promise<WorkOrderAssignment> => {
  const response = await api.patch<ApiResponse<WorkOrderAssignment>>(
    `/work-orders/${workOrderId}/reassign-technician`,
    { technician_id: technicianId },
  );

  if (!response.data.success) {
    throw new Error(
      response.data.message || "Failed to reassign technician",
    );
  }

  return response.data.data;
};


export const assignVehicle = async (
  workOrderId: number,
  vehicleId: number
): Promise<WorkOrderVehicle> => {
  const response = await api.post<ApiResponse<WorkOrderVehicle>>(
    `/work-orders/${workOrderId}/assign-vehicle`,
    {
      vehicle_id: vehicleId,
    }
  );

  if (!response.data.success) {
    throw new Error(
      response.data.message || "Failed to assign vehicle"
    );
  }

  return response.data.data;
};



export const cancelWorkOrder = async (
  workOrderId: number,
): Promise<WorkOrder> => {
  const response = await api.patch<ApiResponse<WorkOrder>>(
    `/work-orders/${workOrderId}/cancel`,
  );

  if (!response.data.success) {
    throw new Error(
      response.data.message || "Failed to cancel work order",
    );
  }

  return response.data.data;
};