export default `
export class GameOfLife {
	private gridSize: number;
	private canvasSize: number;
	private lineColor: string;
	private liveColor: string;
	private deadColor: string;
	private initialLifeProbability: number;
	private animationRate: number;
	private cellSize: number;
	private context: CanvasRenderingContext2D;
	private world;

`;