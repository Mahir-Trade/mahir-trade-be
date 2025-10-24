import { Router, Request, Response } from "express";

// Controllers
import { AuthController } from "../controllers/auth.controller";
import { GroupController } from "../controllers/group.controller";
import { ModuleController } from "../controllers/module.controller";
import { AdminController } from "../controllers/admin.controller";
import { PackageController } from "../controllers/package.controller";
import { SubModuleController } from "../controllers/subModule.controller";
import { ReportController } from "../controllers/report.controller";
import { PaymentController } from "../controllers/payment.controller";
import { SchedulerController } from "../controllers/scheduler.controller";
import { RiskRewardController } from "../controllers/riskReward.controller";
import { SectionController } from "../controllers/section.controller";
import { SectionItemController } from "../controllers/sectionItem.controller";
import { EconomicCalendarController } from "../controllers/economicCalendar.controller";
import { authMiddleware } from "../middlewares/authMiddlewares";

const router = Router();

// Initialize controllers
const authController = new AuthController();
const groupController = new GroupController();
const moduleController = new ModuleController();
const adminController = new AdminController();
const packageController = new PackageController();
const subModuleController = new SubModuleController();
const reportController = new ReportController();
const paymentController = new PaymentController();
const cronController = new SchedulerController();
const riskRewardController = new RiskRewardController();
const sectionController = new SectionController();
const sectionItemController = new SectionItemController();
const economicCalendarController = new EconomicCalendarController();

/* ========================= HEALTH CHECK ========================= */
router.get("/", async (req: Request, res: Response) => {
  res.status(200).json({
    status: 200,
    message: "Service is running!",
  });
});

/* ========================= USERS ========================= */
router.post(`/users/register`, authController.userRegistration);
router.post(`/users/login`, authController.userLogin);
router.get(`/users/login/google`, authController.loginWithGoogle);
router.get(`/users/login/google-callback`, authController.callbackGoogle);
router.post(`/users/forgot-password`, authController.forgotPassword);
router.post(`/users/reset-password`, authController.requestResetPassword);
router.get(`/users/detail`, authController.getDetailUser);

/* ========================= GROUPS ========================= */
router.get(`/groups`, authMiddleware, groupController.getGroups);
router.get(`/groups/:id`, authMiddleware, groupController.getGroupByID);
router.post(`/groups`, authMiddleware, groupController.createGroup);
router.put(`/groups/:id`, authMiddleware, groupController.updateGroup);
router.delete(`/groups/:id`, authMiddleware, groupController.deleteGroup);

/* ========================= DISCORD ========================= */
router.get(`/discord/account`, authController.inviteDiscordUserToGuild);
router.get(
  `/discord/account/add-role`,
  authController.connectDiscordAccountAndAssignRole
);
router.get(
  `/discord/account/remove-role`,
  authController.connectDiscordAccountAndRemoveRole
);
router.post(`/discord/connect-role`, authController.assignRoleDiscordToUser);
router.post(`/discord/remove-role`, authController.removeRoleDiscordUser);

/* ========================= PACKAGES ========================= */
router.get(`/packages`, authMiddleware, packageController.getPackages);
router.get(`/packages/:id`, authMiddleware, packageController.getPackageByID);
router.post(`/packages`, authMiddleware, packageController.createPackage);
router.put(`/packages/:id`, authMiddleware, packageController.updatePackage);
router.delete(`/packages/:id`, authMiddleware, packageController.deletePackage);

/* ========================= MODULES ========================= */
router.get(`/modules`, authMiddleware, moduleController.getModules);
router.get(
  `/modules/:module_id`,
  authMiddleware,
  moduleController.getModuleByID
);
router.get(
  `/modules/group/:group_id`,
  authMiddleware,
  moduleController.getModulesByGroupID
);
router.post(`/modules`, authMiddleware, moduleController.createModule);
router.patch(
  `/modules/:module_id`,
  authMiddleware,
  moduleController.updateModule
);
router.delete(
  `/modules/:module_id`,
  authMiddleware,
  moduleController.deleteModule
);
router.get(
  `/modules/user/:module_id`,
  authMiddleware,
  moduleController.getPercentageMarkWatchedModulesUser
);

/* ========================= SUBMODULES ========================= */
router.get(`/sub-modules`, subModuleController.getSubModules);
router.get(`/sub-modules/:sub_module_id`, subModuleController.getSubModuleByID);
router.get(
  `/sub-modules/module/:module_id`,
  subModuleController.getSubModulesByModuleID
);
router.post(`/sub-modules`, subModuleController.createSubModule);
router.patch(
  `/sub-modules/:sub_module_id`,
  subModuleController.updateSubModule
);
router.delete(
  `/sub-modules/:sub_module_id`,
  subModuleController.softDeleteSubModule
);
router.post(
  `/sub-modules/mark-watched`,
  subModuleController.markSubModuleAsWatched
);

/* ========================= ADMINS ========================= */
router.post(`/admins/register`, adminController.adminRegistration);
router.post(`/admins/login`, adminController.adminLogin);
router.get(
  `/admins/detail`,
  authMiddleware,
  adminController.getDetailAdminInfo
);
router.get(`/admins/user-detail/:user_id`, adminController.getDetailUserForBO);
router.get(`/admins/users`, adminController.getAllUsers);

router.post(
  "/admins/users/toggle-expired",
  adminController.toggleInactiveUserMembership
);
router.post(
  "/admins/start-membership-program",
  authMiddleware,
  adminController.startMembershipProgram
);
router.get(
  "/admins/membership-program",
  authMiddleware,
  adminController.getMembershipProgramDate
);
router.put(
  "/admins/membership-program",
  authMiddleware,
  adminController.updateMembershipProgramDate
);
//

/* ========================= REPORTS ========================= */
router.get(`/reports`, reportController.getReports);
router.get(`/reports/:id`, reportController.getReportByID);
router.post(`/reports`, reportController.createReport);
router.put(`/reports/:id`, reportController.updateReport);
router.delete(`/reports/:id`, reportController.deleteReport);

/* ========================= ORDERS ========================= */
router.post(`/orders/create`, paymentController.createPayment);

/* ========================= UPLOAD ========================= */
router.post(`/upload`, subModuleController.uploadFile);

/* ========================= PUBLIC ========================= */
router.post(`/payment-link-callback`, paymentController.paymentLinkCallback);
router.get(`/cron`, paymentController.paymentLinkCallback);
router.get(`/package-list`, packageController.getPackages);

/* ========================= SECTION ========================= */
router.post(`/section/create`, sectionController.createSection);
router.patch(`/section/:id`, sectionController.updateSection);
router.get(`/section/getAll`, sectionController.getFullSection);
router.delete(`/section/:id`, sectionController.deleteSection);
router.get(`/section/:type`, sectionController.getSectionByType);

/* ========================= SECTION ITEMS ========================= */
router.get(`/section-item/getAll`, sectionItemController.getItems);
router.post(`/section-item/create`, sectionItemController.createItem);
router.patch(`/section-item/:id`, sectionItemController.updateItem);
router.get(
  `/section-item/:section_id`,
  sectionItemController.getItemsBySectionID
);
router.delete(`/section-item/:id`, sectionItemController.deleteItem);

/* ========================= ECONOMIC CALENDAR ========================= */
router.get(`/economic-calendar/getAll`, economicCalendarController.getCalendar);

export default router;
