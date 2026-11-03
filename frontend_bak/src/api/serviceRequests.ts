
import { api } from "./client";


export interface ServiceRequest {
  id: number;
  customer_id: number;
  title: string;
  description: string;
  service_address: string;
  preferred_date: string | null;
  status: string;
  created_at: string;
  updated_at: string;
}


interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

function getErrorMessage(error: unknown, fallback: string): string {
  if (error instanceof Error && error.message) {
    return error.message;
  }

  return fallback;
}

export const serviceRequestApi = {
  async getAll(): Promise<ServiceRequest[]> {
    try {
      const response = await api.get<ApiResponse<ServiceRequest[]>>(
        "/service-requests",
      );

      if (!response.data.success) {
        throw new Error(
          response.data.message || "Failed to load service requests",
        );
      }

      return response.data.data;
    } catch (error) {
      throw new Error(
        getErrorMessage(error, "Failed to load service requests"),
      );
    }
  },

  async getById(id: number): Promise<ServiceRequest> {
    try {
      const response = await api.get<ApiResponse<ServiceRequest>>(
        `/service-requests/${id}`,
      );

      if (!response.data.success) {
        throw new Error(
          response.data.message || "Failed to load service request",
        );
      }

      return response.data.data;
    } catch (error) {
      throw new Error(
        getErrorMessage(error, "Failed to load service request"),
      );
    }
  },

    async create(payload: {
        title: string;
        description: string;
        service_address: string;
        preferred_date: string | null;
    }): 
    Promise<ServiceRequest> {
        const response = await api.post<ApiResponse<ServiceRequest>>(
            "/service-requests",
            payload
        );

        if (!response.data.success) {
            throw new Error(
            response.data.message || "Failed to create service request"
            );
        }

        return response.data.data;
    },
};
