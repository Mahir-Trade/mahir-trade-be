import cron from "node-cron";
import { UserService } from "../services/user.service";

export class SchedulerController {
  private userService: UserService;

  constructor() {
    this.userService = new UserService();

    this.registerJobs();
  }

  private registerJobs(): void {
    const jobs = [
      {
        spec: "59 23 * * *", // sama persis dengan cron Go
        cmd: async () => {
          try {
            await this.userService.updateMembership();
          } catch (err: any) {
            console.error(
              "[SchedulerController] failed to update user status:",
              err
            );
          }
        },
      },
    ];

    jobs.forEach((job) => {
      cron.schedule(job.spec, job.cmd);
    });
  }
}
