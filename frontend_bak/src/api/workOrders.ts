
import { api } from "./client";

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
  created_at: string;
  updated_at: string;
}

interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

export const workOrderApi = {
  async getAll(): Promise<WorkOrder[]> {
    const response = await api.get<ApiResponse<WorkOrder[]>>(
      "/work-orders",
    );

    if (!response.data.success) {
      throw new Error(response.data.message || "Failed to load work orders");
    }

    return response.data.data;
  },
};
